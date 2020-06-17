import React, {Component} from "react";
import ServiceList from "./ServiceList";
import ActivatedList from "./ActivatedList";
import SituationList from "./SituationList"
import CreatorBox from "./CreatorBox";
import Loading from "./Loading";
import Error from "./Error";
import "./SituationServices.css"


const situationIndex = 0;
const serviceIndex = 1;
const activatedServiceIndex = 2;

const SERVICE = 0;
const SITUATION = 1;

const urls = [
    "http://moodybase:8080/situation/",
    "http://moodybase:8080/service/",
    "http://moodybase:7000/sensor_service",
]

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
        .catch(error => this.setState({isLoaded: true, error}))
    }

    handleSituationRemoval(id) {
        const situationList = this.state.situationList;

        fetch(urls[situationIndex] + id, {method: "DELETE"})
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

        fetch(urls[serviceIndex] + id, {method: "DELETE"})
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
                })}
            );
    }

    handleServiceActivation(id) {
        const activateRequest = {serviceId: id};
        const jsonRequest = JSON.stringify(activateRequest);

        const services = this.state.serviceList;
        let activatedServiceList = this.state.activatedServiceList;

        fetch(urls[activatedServiceIndex], {method: "POST", body: jsonRequest})
        .then(response => response.json())
        .then(
            result => {
                activatedServiceList.push(result.name);
                this.setState({
                    isLoaded: true,
                    serviceList: services.filter(service => service.id !== id),
                    activatedServiceList: activatedServiceList
                });
            },
            error => this.setState({
                isLoaded: true,
                error
            })
        )
    }

    handleServiceDeactivation(name) {
        const activateRequest = {name: name};
        const jsonRequest = JSON.stringify(activateRequest);

        let serviceList = this.state.serviceList;
        const activatedServiceList = this.state.activatedServiceList;
        const allServices = this.state.allServices;

        fetch(urls[activatedServiceIndex], {method: "DELETE", body: jsonRequest})
        .then(resp => resp.json())
        .then(
            result => {
                const targetServiceToAdd = allServices.find(service => service.name === result.name)
                serviceList.push(targetServiceToAdd);
                this.setState({
                    isLoaded: true,
                    serviceList: serviceList,
                    activatedServiceList: activatedServiceList.filter(service => service !== result.name)
                })},
            error => this.setState({
                isLoaded: true,
                error
            })
        )
    }

    //Entity creation

    createNew(entity, name) {
        const targetUrl = entity === SITUATION ? situationIndex : serviceIndex;
        const createRequest = {name: name};
        const jsonRequest = JSON.stringify(createRequest);

        let dataList = entity === SITUATION ? this.state.situationList : this.state.serviceList;

        fetch(urls[targetUrl], {method: "POST", body: jsonRequest})
            .then(resp => resp.json())
            .then(
                response => {
                    dataList.push(response);
                    if(entity === SITUATION){
                        this.setState({isLoaded: true, situationList: dataList})
                    } else {
                        let allServices = this.state.allServices;
                        allServices.push(response);
                        this.setState({isLoaded: true, serviceList: dataList, allServices: allServices})
                    }
                },
                error => this.setState({isLoaded: true, error})
            );
    }

    handleCreateSituation(name) {
        this.createNew(SITUATION, name);
    }

    handleCreateService(name) {
        this.createNew(SERVICE, name);
    }



    render() {
        const {situationList, serviceList, activatedServiceList, isLoaded, error} = this.state;
        if(error) {
            console.log(error);
            return <Error name={error} />;
        }

        if(!isLoaded) {
            return <Loading />;
        } else {
            return(
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
}