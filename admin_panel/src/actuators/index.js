import React, { Component } from "react";
import socketIOClient from "socket.io-client";
import { Empty, Error, Loading } from "./Errors"
import { ActivateServerMode, DeactivateServerMode } from "./ActuatorModeBoxes"
import ActuatorGrid from "./ActuatorGrid";
import "./Boxes.css";

import Actuator from "./models/Actuator"
import IpMappingReq from "./models/IpMappingReq"
import SwitchModeReq from "./models/SwitchModeReq"
import IpToDelete from "./models/IpToDelete"

const activatedActuatorServerIndex = 0;
const situationIndex = 1;
const activateModeIndex = 2;

const socketioEndpoint = "http://localhost:7000";
const socketioEvent = "actuator";

const urls = [
    "http://localhost:7000/actuators",
    "http://localhost:8080/situation/",
    "http://localhost:7000/actuator_mode",
];


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
    }

    componentDidMount() {
        const fetchPromises = urls.map(url => fetch(url).then(response => response.json()))
        Promise.all(fetchPromises)
            .then(responses =>
                this.setState({
                    actuatorMode: responses[activateModeIndex].mode,
                    actuatorList: responses[activatedActuatorServerIndex].actuatorList.map(actuator => new Actuator(actuator.ip, actuator.mappings)),
                    situationList: responses[situationIndex].situations,
                    isLoaded: true,
                })
            )
            .catch(
                error => this.setState({ isLoaded: true, error })
            );

        const socket = socketIOClient(socketioEndpoint);
        socket.on(socketioEvent, data => {
            let rawActuator = JSON.parse(atob(data));
            let dataActuator = new Actuator(rawActuator.ip, rawActuator.mappings);
            let actuators = this.state.actuatorList;
            actuators.push(dataActuator);
            this.setState({ isLoaded: true, actuatorList: actuators });
        });
    }

    addMapping(ip, situation, action) {
        if (!ip || !situation || !action) {
            alert("Can't have any empty field!");
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
                    actuators[targetIndex].mappingList.push(response);
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
        if (!isLoaded) return <div className="actuators"><Loading /></div>
        if (actuatorMode) return <div className="actuators"><ActivateServerMode handleSwitch={() => this.switchActuatorMode(false)} /></div>
        if (actuatorList.length === 0) return <div className="actuators"><div className="empty"><Empty /><DeactivateServerMode handleSwitch={() => this.switchActuatorMode(true)} /></div></div>;

        return (
            <div className="actuators">
                <ActuatorGrid handleSwitchMode={() => this.switchActuatorMode(true)} handleAdd={this.addMapping} handleDelete={this.deleteMappings} actuatorList={actuatorList} situationList={situationList} />
            </div>
        );
    };
}