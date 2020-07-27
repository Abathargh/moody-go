import React, { Component } from "react";

export default class Mapping extends Component {
    render() {
        return (
            <tr>
                <td><span>{this.props.situation}</span></td>
                <td><span>{this.props.action}</span></td>
            </tr>
        );
    }
}
