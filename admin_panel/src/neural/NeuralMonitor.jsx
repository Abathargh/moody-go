import React, { Component } from "react";
import ChangeSettings from "./ChangeSettings";
import ActiveStateTable from "./ActiveStateTable"

export default class NeuralMonitor extends Component {
    render() {
        return (
            <div className="neuralBox ">
                <h2>Neural Monitor</h2>
                <ActiveStateTable neuralState={this.props.neuralState} />
                <ChangeSettings neuralState={this.props.neuralState} handleStateChange={this.props.handleStateChange} datasetList={this.props.datasetList} />
            </div >
        );
    }
}