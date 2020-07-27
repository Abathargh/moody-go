import React, { Component } from "react"

class ActivateServerMode extends Component {
    render() {
        return (
            <div className="activateModeBox">
                <button onClick={this.props.handleSwitch}>Activate Server Mode</button>
            </div>
        );
    }
}

class DeactivateServerMode extends Component {
    render() {
        return (
            <div className="activateModeBox actuatorBox">
                <button onClick={this.props.handleSwitch}>End Server Mode</button>
            </div>
        );
    }
}

export { ActivateServerMode, DeactivateServerMode };