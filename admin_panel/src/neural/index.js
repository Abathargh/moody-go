import React, { Component } from "react";
import DatasetList from "./DatasetList";
import NeuralMonitor from "./NeuralMonitor";
import SituationBox from "./SituationBox";
import NewDatasetBox from "./NewDatasetBox";
import { Error, Loading } from "./Errors";

import NeuralState, { Mode } from "./models/NeuralState"
import SituationState from "./models/SituationState"
import SituationRequest from "./models/SituationRequest"
import Dataset from "./models/DataSet"

import "./NeuralStyle.css"


const datasetIndex = 0;
const gatewayIndex = 1;
const situationIndex = 2;
const getSituationIndex = 3;
const sensorServiceIndex = 4;

let url = new URL(window.location.origin);
url.port = process.env.REACT_APP_API_PORT;
let gateway = url.origin;

const urls = [
    "/dataset",
    "/neural_state",
    "/current_situation",
    "/situation",
    "/sensor_service",
].map(url => gateway + url);

export default class Neural extends Component {
    constructor(props) {
        super(props);

        this.state = {
            currentSituation: "-",
            situationList: [],
            activeServices: [],
            neuralState: null,
            datasetList: [],
            isLoaded: false,
            error: null
        }

        this.handleStateChange = this.handleStateChange.bind(this);
        this.handleSituationChange = this.handleSituationChange.bind(this);
        this.handleDatasetRemove = this.handleDatasetRemove.bind(this);
        this.handleDatasetCreate = this.handleDatasetCreate.bind(this);
    }

    componentDidMount() {
        const fetchPromises = urls.map(url => fetch(url).then(resp => resp.json()));
        Promise.all(fetchPromises)
            .then(response =>
                this.setState({
                    currentSituation: SituationState.FromAPI(response[situationIndex]),
                    situationList: response[getSituationIndex].situations,
                    activeServices: response[sensorServiceIndex].services,
                    neuralState: NeuralState.FromAPI(response[gatewayIndex]),
                    datasetList: response[datasetIndex].datasets,
                    isLoaded: true
                })
            )
            .catch(error => this.setState({ isLoaded: true, error: error }))
    }

    handleStateChange(newMode, newDataset) {
        console.log(newMode, newDataset);
        if (!newMode || (!newDataset && newMode !== Mode.Stopped)) {
            alert("Can't have an empty mode/dataset field!")
            return;
        }

        const datasets = this.state.datasetList;
        if (newMode !== Mode.Stopped && !datasets.some(dataset => dataset.name === newDataset)) {
            alert("There's no dataset with the passed name!");
            return;
        }

        const oldState = this.state.neuralState;
        const newState = new NeuralState(newMode, newDataset);
        if (!newState.equals(oldState)) {
            const stateReq = JSON.stringify(newState.toReqFormat());
            fetch(urls[gatewayIndex], { method: "PUT", body: stateReq })
                .then(resp => resp.json())
                .then(
                    response => this.setState({ neuralState: NeuralState.FromAPI(response), isLoaded: true }),
                    error => this.setState({ isLoaded: true, error: error })
                )
        }
    }

    handleSituationChange(situationId) {
        if (!situationId) {
            alert("Select a situation to be set!");
            return
        }

        const currentSituation = this.state.currentSituation;
        if (currentSituation.sameId(situationId)) {
            return;
        }

        var situationReq = JSON.stringify(new SituationRequest(situationId));
        fetch(urls[situationIndex], { method: "PUT", body: situationReq })
            .then(resp => resp.json())
            .then(
                response => this.setState({ isLoaded: true, currentSituation: SituationState.FromAPI(response) }),
                error => this.setState({ isLoaded: true, error: error })
            )
    }

    handleDatasetRemove(name) {
        if (!name) {
            alert("There was an error in removing the dataset!");
            return;
        }

        const datasets = this.state.datasetList;
        if (!datasets.some(dataset => dataset.name === name)) {
            alert("No such dataset!");
            return;
        }

        fetch(urls[datasetIndex] + "/" + name, { method: "DELETE" })
            .then(resp => resp.json())
            .then(
                response => {
                    const datasets = this.state.datasetList;
                    this.setState({ isLoaded: true, datasetList: datasets.filter(dataset => dataset.name !== response.name) });
                },
                error => this.setState({ isLoaded: true, error: error })
            )
    }


    handleDatasetCreate(name) {
        if (!name) {
            alert("Insert a dataset name!");
            return;
        }

        const datasets = this.state.datasetList;
        if (datasets.some(dataset => dataset.name === name)) {
            alert("A dataset with this name already exists!");
            return;
        }

        const activeServices = this.state.activeServices;
        if (activeServices.length === 0) {
            alert("You must have at least one active service to create a new dataset!");
            return;
        }

        const keys = this.state.activeServices;
        const newDataset = JSON.stringify(new Dataset(name, keys));
        fetch(urls[datasetIndex], {
            headers: {
                "Content-Type": "application/json"
            },
            method: "POST",
            body: newDataset
        })
            .then(resp => resp.json())
            .then(
                response => {
                    const datasets = this.state.datasetList;
                    datasets.push(response);
                    this.setState({ isLoaded: true, datasetList: datasets })
                },
                error => this.setState({ isLoaded: true, error: error })
            )
    }

    render() {
        const { currentSituation, situationList, neuralState, datasetList, isLoaded, error } = this.state;
        console.log(currentSituation);
        if (error) return <div className="neuralMonitor"><Error error={error.toString()} /></div>
        if (!isLoaded) return <div className="neuralMonitor"><Loading /></div>
        return (
            <div className="neuralMonitor">
                <NeuralMonitor handleStateChange={this.handleStateChange} neuralState={neuralState} datasetList={datasetList} />
                <SituationBox changeSituation={this.handleSituationChange} currentSituation={currentSituation} situationList={situationList} />
                <DatasetList handleRemove={this.handleDatasetRemove} datasetList={datasetList} />
                <NewDatasetBox handleCreate={this.handleDatasetCreate} />
            </div>
        );
    };
}