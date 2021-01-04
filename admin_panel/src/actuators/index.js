import { ActivateServerMode, DeactivateServerMode } from "./ActuatorModeBoxes"
import { w3cwebsocket as W3CWebSocket } from "websocket";
import { Empty, Error, Loading } from "./Errors"
import ActuatorGrid from "./ActuatorGrid";
import React, { Component } from "react";
import "./Boxes.css";

import Actuator from "./models/Actuator"
import IpMappingReq from "./models/IpMappingReq"
import SwitchModeReq from "./models/SwitchModeReq"
import IpToDelete from "./models/IpToDelete"

const activatedActuatorServerIndex = 0;
const situationIndex = 1;
const activateModeIndex = 2;

let url = new URL(window.location.origin);
url.port = process.env.REACT_APP_API_PORT;
let gateway = url.origin;

const urls = [
    "/actuators",
    "/situation",
    "/actuator_mode",
].map(url => gateway + url);

url.protocol = "ws:";
const wsURL = url.origin + "/actuator_ws";

const client = new W3CWebSocket(wsURL);
client.onopen = () => {
    console.log('WebSocket Client Connected');
};



export default class Actuators extends Component {
    constructor(props) {
        super(props);
        this.state = {
            actuatorMode: true,
            actuatorList: [],
            situationList: [],
            isLoaded: false,
            error: null,
        }

        this.addMapping = this.addMapping.bind(this);
        this.deleteMappings = this.deleteMappings.bind(this);
        this.switchActuatorMode = this.switchActuatorMode.bind(this);

        client.onmessage = (message) => {
            //atob
            let rawActuator = JSON.parse(message.data);
            if(Actuator.isActuatorData(rawActuator)) {
                const situations = this.state.situationList;
                let dataActuator = Actuator.init(rawActuator.ip, rawActuator.mappings, situations);
                let actuators = this.state.actuatorList;
                actuators.push(dataActuator);
                this.setState({ isLoaded: true, actuatorList: actuators });
            }
        };
    }

    componentDidMount() {
        const fetchPromises = urls.map(url => fetch(url).then(response => response.json()))
        Promise.all(fetchPromises)
            .then(responses => {
                const situations = responses[situationIndex].situations;
                const actuators = responses[activatedActuatorServerIndex].actuatorList;
                this.setState({
                    actuatorMode: responses[activateModeIndex].mode,
                    actuatorList: actuators.map(actuator => Actuator.init(actuator.ip, actuator.mappings, situations)),
                    situationList: situations,
                    isLoaded: true,
                })
            })
            .catch(
                error => this.setState({ isLoaded: true, error })
            );

    }

    addMapping(ip, situation, action) {
        if (!ip || !situation || !action) {
            alert("Can't have any empty field!");
            return
        }

        action = Number(action);
        if (!Number.isInteger(action) || action < 0 || action > 255) {
            alert("Actions must be codified with a byte (0 < x < 255)");
            return
        }

        let actuators = this.state.actuatorList;
        if (!actuators.some(actuator => actuator.ip === ip)) {
            alert("There's no actuator with such IP!");
            return
        }

        let targetIndex = actuators.findIndex(actuator => actuator.ip === ip);
        if (actuators[targetIndex].mappingList.some(mapping => mapping.situation === parseInt(situation))) {
            alert("A mapping with this situation already exists!");
            return
        }

        const ipMappingReq = JSON.stringify(new IpMappingReq(ip, situation, action));
        fetch(urls[activatedActuatorServerIndex], { method: "POST", body: ipMappingReq })
            .then(resp => resp.json())
            .then(
                response => {
                    const situations = this.state.situationList;
                    actuators[targetIndex].mappingList.push(Actuator.mappingWithName(response, situations));
                    this.setState({ actuatorList: actuators, isLoaded: true })
                },
                error => this.setState({ isLoaded: true, error: error })
            )
    }

    deleteMappings(ip) {
        if (!ip) {
            alert("Can't have an empty IP field!");
            return
        }

        let actuators = this.state.actuatorList;
        if (!actuators.some(actuator => actuator.ip === ip)) {
            alert("There's no actuator with such IP!");
            return
        }

        if (actuators.length !== 0) {
            const ipToDelete = JSON.stringify(new IpToDelete(ip));
            fetch(urls[activatedActuatorServerIndex], { method: "DELETE", body: ipToDelete })
                .then(
                    _ => {
                        const actuators = this.state.actuatorList;
                        const targetIndex = actuators.findIndex(actuator => actuator.ip === ip);
                        actuators[targetIndex].mappingList = [];
                        this.setState({ actuatorList: actuators, isLoaded: true })
                    },
                    error => this.setState({ isLoaded: true, error: error })
                )
        }
    }

    switchActuatorMode(state) {
        const mode = JSON.stringify(new SwitchModeReq(state));
        fetch(urls[activateModeIndex], { method: "POST", body: mode })
            .then(resp => resp.json())
            .then(
                response => this.setState({ actuatorMode: response.mode, actuatorList: [], isLoaded: true }),
                error => this.setState({ isLoaded: true, error: error })
            )
    }


    render() {
        const { actuatorMode, actuatorList, situationList, isLoaded, error } = this.state;
        if (error) return <div className="actuators"><Error error={error.toString()} /></div>
        if (!isLoaded) return <div className="actuators"><div className="empty"><Loading /><DeactivateServerMode handleSwitch={() => this.switchActuatorMode(true)} /></div></div>
        if (actuatorMode) return <div className="actuators"><ActivateServerMode handleSwitch={() => this.switchActuatorMode(false)} /></div>
        if (actuatorList.length === 0) return <div className="actuators"><div className="empty"><Empty /><DeactivateServerMode handleSwitch={() => this.switchActuatorMode(true)} /></div></div>;

        return (
            <div className="actuators">
                <ActuatorGrid handleSwitchMode={() => this.switchActuatorMode(true)} handleAdd={this.addMapping} handleDelete={this.deleteMappings} actuatorList={actuatorList} situationList={situationList} />
            </div>
        );
    };
}