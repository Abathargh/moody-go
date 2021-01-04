import { w3cwebsocket as W3CWebSocket } from "websocket";
import { Empty, Error, Loading } from "./Errors"
import ServiceData from "./models/ServiceData";
import React, { Component } from "react";
import ServiceGrid from "./ServiceGrid";

let url = new URL(window.location.origin);

url.port = process.env.REACT_APP_API_PORT
const serviceURL = url.origin + "/data_table";

url.protocol = "ws:";
const wsURL = url.origin + "/service_ws";

const client = new W3CWebSocket(wsURL);
client.onopen = () => {
    console.log('WebSocket Client Connected');
};

export default class Monitor extends Component {
    constructor(props) {
        super(props);
        this.state = {
            serviceList: [],
            isLoaded: false,
            error: null,
        }

        client.onmessage = (message) => {
            const jsonData = JSON.parse(message.data);
            if(ServiceData.isServiceData(jsonData)) {
                let dataService = new ServiceData(jsonData.service, jsonData.data);
                let services = this.state.serviceList;
                let index = services.findIndex(service => service.service === dataService.service);
                if (index === -1) { return; }
                services[index].data = dataService.data;
                this.setState({ isLoaded: true, serviceList: services });
            }
        };
    }

    async componentDidMount() {
        const resp = await fetch(serviceURL);
        try {
            const respObj = await resp.json();
            let serviceList = [];
            for(let [key, value] of Object.entries(respObj.services)) {
                serviceList.push(new ServiceData(key, value));
            }
            this.setState({ isLoaded: true, serviceList: serviceList })
        } catch(err) {
            this.setState({ isLoaded: true, error: err })
        }
    }

    render() {
        const { serviceList, isLoaded, error } = this.state;
        if (error) return <div className="monitor"><Error /></div>
        if (!isLoaded) return <div className="monitor"><Loading /></div>
        if (serviceList.length === 0) return <div className="monitor"><Empty /></div>;

        return <div className="monitor"><ServiceGrid serviceList={serviceList} /></div>;
    };
}