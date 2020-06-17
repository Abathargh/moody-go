import React, { Component } from "react";
import {Remover} from "./Activators";

export default class Situation extends Component {
    render() {
        return(
            <tr>
                <td><span>{ this.props.name }</span></td>
                <td><Remover handleRemove={() => this.props.handleRemove(this.props.id)} /></td>
            </tr>
        );
    }
}
