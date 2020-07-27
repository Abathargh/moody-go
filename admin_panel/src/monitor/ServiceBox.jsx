import React, { Component } from "react";

export default class ServiceBox extends Component {
    render() {
        return (
            <div className="serviceBox">
                <h3>{this.props.service}</h3>
                <h3>{this.props.data}</h3>
            </div>
        );
    }
}