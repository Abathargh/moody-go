import React, { Component } from "react";
import DatasetList from "./DatasetList";
import { Empty, Loading } from "./Errors";

import "./NeuralStyle.css"

export default class Neural extends Component {
    constructor(props) {
        super(props);

        this.state = {
            datasetList: [],
            isLoaded: false,
            error: null
        }
    }

    componentDidMount() {
        fetch("http://localhost:8080/dataset")
            .then(resp => resp.json())
            .then(
                response => this.setState({ datasetList: response.datasets, isLoaded: true }),
                error => this.setState({ isLoaded: true, error: error })
            )
    }

    render() {
        const { datasetList, isLoaded, error } = this.state;
        console.log(datasetList);
        if (error) return <div className="neuralMonitor"><Empty /></div>
        if (!isLoaded) return <div className="neuralMonitor"><Loading /></div>
        return <div className="neuralMonitor"><DatasetList datasetList={datasetList} /></div>;
    };
}