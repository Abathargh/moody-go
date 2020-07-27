import React, { Component } from "react";

class Error extends Component {
    render() {
        return <div className="error"><h2>{this.props.error}</h2></div>;
    }
}

class Loading extends Component {
    render() {
        return <div className="error"><h2>Loading...</h2></div>;
    }
}

class Empty extends Component {
    render() {
        return <div className="error actuatorBox"><h2>No active actuator server!</h2></div>;
    }
}

export { Empty, Error, Loading };