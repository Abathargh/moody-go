import React, { Component } from "react";
import socketIOClient from "socket.io-client";
import ServiceGrid from "./ServiceGrid";
import { Empty, Error, Loading } from "./Errors"

const socketioIndex = 0;
const activatedServiceIndex = 1;

const socketioEvent = "data";

let config = require("../conf.json");

const urls = [
    "/",
    "/sensor_service",
].map(url => config.gateway + url);

class Service {
    constructor(obj) {
        this.service = obj.service;
        this.data = obj.data;
    }
}

export default class Monitor extends Component {
    constructor(props) {
        super(props);
        this.state = {
            serviceList: [],
            isLoaded: false,
            error: null,
        }
    }

    componentDidMount() {
        fetch(urls[activatedServiceIndex])
            .then(resp => resp.json())
            .then(
                response => {
                    const completeServices = response.services.map(service => new Service({ service: service, data: "-" }))
                    console.log(completeServices)
                    this.setState({ isLoaded: true, serviceList: completeServices })
                },
                error => this.setState({ isLoaded: true, error })
            )

        const socket = socketIOClient(urls[socketioIndex]);
        socket.on(socketioEvent, data => {
            let dataService = new Service(JSON.parse(atob(data)));
            let services = this.state.serviceList;
            let index = services.findIndex(service => service.service === dataService.service);
            if (index === -1) { return; }
            services[index].data = dataService.data;
            this.setState({ isLoaded: true, serviceList: services });
        });
    }

    render() {
        const { serviceList, isLoaded, error } = this.state;
        if (error) return <div className="monitor"><Error /></div>
        if (!isLoaded) return <div className="monitor"><Loading /></div>
        if (serviceList.length === 0) return <div className="monitor"><Empty /></div>;

        return <div className="monitor"><ServiceGrid serviceList={serviceList} /></div>;
    };
}