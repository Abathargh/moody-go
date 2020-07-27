import React, { Component } from "react";
import ActuatorBox from "./ActuatorBox";
import { DeactivateServerMode } from "./ActuatorModeBoxes"

export default class ActuatorGrid extends Component {
    constructor(props) {
        super(props);

        this.forwardAdd = this.forwardAdd.bind(this);
        this.forwardDelete = this.forwardDelete.bind(this);
    }

    forwardAdd(ip, situation, action) {
        this.props.handleAdd(ip, situation, action)
    }

    forwardDelete(ip) {
        this.props.handleDelete(ip)
    }

    render() {
        return (
            <div className="actuatorGrid">
                {
                    this.props.actuatorList.map(
                        actuator => <ActuatorBox handleAdd={this.forwardAdd} handleDelete={this.forwardDelete} key={actuator.ip} ip={actuator.ip} mappingList={actuator.mappingList} situationList={this.props.situationList} />
                    )
                }
                <DeactivateServerMode handleSwitch={this.props.handleSwitchMode} />
            </div>
        );
    }
}