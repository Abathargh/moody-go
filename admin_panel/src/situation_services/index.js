import React, { Component } from "react";
import ServiceList from "./ServiceList";
import ActivatedList from "./ActivatedList";
import SituationList from "./SituationList"
import CreatorBox from "./CreatorBox";
import Loading from "./Loading";
import "./SituationServices.css"


const situationIndex = 0;
const serviceIndex = 1;
const activatedServiceIndex = 2;

const SERVICE = 0;
const SITUATION = 1;

let url = new URL(window.location.origin);
url.port = process.env.REACT_APP_API_PORT;
let gateway = url.origin;

const urls = [
    "/situation",
    "/service",
    "/sensor_service",
].map(url => gateway + url);

export default class ServiceSituations extends Component {
    situations;
    constructor(props) {
        super(props);
        this.state = {
            allServices: [],
            situationList: [],
            serviceList: [],
            activatedServiceList: [],
            isLoaded: false,
            error: null
        };
        this.handleSituationRemoval = this.handleSituationRemoval.bind(this);
        this.handleServiceActivation = this.handleServiceActivation.bind(this);
        this.handleServiceDeactivation = this.handleServiceDeactivation.bind(this);
        this.handleServiceRemoval = this.handleServiceRemoval.bind(this);
        this.handleCreateSituation = this.handleCreateSituation.bind(this);
        this.handleCreateService = this.handleCreateService.bind(this);
    }


    componentDidMount() {
        const fetchPromises = urls.map(url => fetch(url).then(response => response.json()))
        Promise.all(fetchPromises)
            .then(responses => this.setState({
                isLoaded: true,
                allServices: responses[serviceIndex].services,
                situationList: responses[situationIndex].situations,
                serviceList: responses[serviceIndex].services.filter(service => !responses[activatedServiceIndex].services.includes(service.name)),
                activatedServiceList: responses[activatedServiceIndex].services
            }), error => console.log(error))
            .catch(error => this.setState({ isLoaded: true, error }));
    }

    handleSituationRemoval(id) {
        const situationList = this.state.situationList;

        fetch(urls[situationIndex] + "/" + id, { method: "DELETE" })
            .then(resp => resp.json())
            .then(
                result => {
                    let removedSituations = situationList.filter(s => s.id !== result.id)
                    this.setState({
                        isLoaded: true,
                        situationList: removedSituations
                    })
                },
                error => this.setState({
                    isLoaded: true,
                    error
                })
            );
    }

    handleServiceRemoval(id) {
        const serviceList = this.state.serviceList;
        const allServices = this.state.allServices;

        fetch(urls[serviceIndex] + "/" + id, { method: "DELETE" })
            .then(resp => resp.json())
            .then(
                result => {
                    let removedServices = serviceList.filter(s => s.id !== result.id);
                    let removedAllServices = allServices.filter(s => s.id !== result.id);
                    this.setState({
                        isLoaded: true,
                        serviceList: removedServices,
                        allServices: removedAllServices,
                    });
                },
                error => {
                    this.setState({
                        isLoaded: true,
                        error
                    })
                }
            );
    }

    handleServiceActivation(id) {
        const activateRequest = { serviceId: id };
        const jsonRequest = JSON.stringify(activateRequest);

        const services = this.state.serviceList;
        let activatedServiceList = this.state.activatedServiceList;

        fetch(urls[activatedServiceIndex], { method: "POST", body: jsonRequest })
            .then(response => response.json())
            .then(
                result => {
                    activatedServiceList.push(result.name);
                    this.setState({
                        isLoaded: true,
                        serviceList: services.filter(service => service.id !== id),
                        activatedServiceList: activatedServiceList
                    });
                    if(result.name === "kefka") console.log("4c6966652e2e2e20647265616d732e2e2e20686f70652e2e2e20576865726520646f207468657920636f6d652066726f6d3f20416e6420776865726520646f207468657920676f3f2053756368206d65616e696e676c657373207468696e67732e2e2e49276c6c2064657374726f79207468656d20616c6c21")
                },
                error => this.setState({
                    isLoaded: true,
                    error: error
                })
            )
    }

    handleServiceDeactivation(name) {
        const activateRequest = { name: name };
        const jsonRequest = JSON.stringify(activateRequest);

        let serviceList = this.state.serviceList;
        const activatedServiceList = this.state.activatedServiceList;
        const allServices = this.state.allServices;

        fetch(urls[activatedServiceIndex], { method: "DELETE", body: jsonRequest })
            .then(resp => resp.json())
            .then(
                result => {
                    const targetServiceToAdd = allServices.find(service => service.name === result.name)
                    serviceList.push(targetServiceToAdd);
                    this.setState({
                        isLoaded: true,
                        serviceList: serviceList,
                        activatedServiceList: activatedServiceList.filter(service => service !== result.name)
                    })
                },
                error => this.setState({
                    isLoaded: true,
                    error: error
                })
            )
    }

    //Entity creation

    createNew(entity, name) {
        const targetUrl = entity === SITUATION ? situationIndex : serviceIndex;
        const createRequest = { name: name };
        const jsonRequest = JSON.stringify(createRequest);

        let dataList = entity === SITUATION ? this.state.situationList : this.state.serviceList;

        fetch(urls[targetUrl], { method: "POST", body: jsonRequest })
            .then(resp => resp.json())
            .then(
                response => {
                    dataList.push(response);
                    if (entity === SITUATION) {
                        this.setState({ isLoaded: true, situationList: dataList })
                    } else {
                        let allServices = this.state.allServices;
                        allServices.push(response);
                        this.setState({ isLoaded: true, serviceList: dataList, allServices: allServices })
                    }
                },
                error => this.setState({ isLoaded: true, error: error })
            );
    }

    handleCreateSituation(name) {
        const situations = this.state.situationList;
        if (!situations.some(situation => situation.name === name)) {
            this.createNew(SITUATION, name);
            return
        }
        alert("The situation alredy exists!")
    }

    handleCreateService(name) {
        const services = this.state.serviceList;
        if (!services.some(service => service.name === name)) {
            this.createNew(SERVICE, name);
            return
        }
        alert("The service alredy exists!")
    }



    render() {
        const { situationList, serviceList, activatedServiceList, isLoaded, error } = this.state;
        if (error) {
            console.log(error);
        }

        if (!isLoaded) {
            return <Loading />;
        }

        return (
            <div className="service_situations">
                <div className="column">
                    <ServiceList
                        serviceList={serviceList}
                        handleServiceActivation={this.handleServiceActivation}
                        handleServiceRemoval={this.handleServiceRemoval}
                    />
                    <ActivatedList
                        activatedServiceList={activatedServiceList}
                        handleServiceDeactivation={this.handleServiceDeactivation}
                    />
                </div>
                <div className="column">
                    <SituationList
                        situationList={situationList}
                        handleSituationRemoval={this.handleSituationRemoval}
                    />
                </div>
                <div className="column">
                    <CreatorBox
                        name="Situation"
                        handleCreate={this.handleCreateSituation}
                    />
                    <CreatorBox
                        name="Service"
                        handleCreate={this.handleCreateService}
                    />
                </div>
            </div>
        );
    }
}

