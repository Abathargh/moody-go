import React, { Component } from "react";
import {Activator, Deactivator, Remover} from "./Activators";

class Service extends Component {
    render() {
        return(
            <tr>
                <td><span>{ this.props.name }</span></td>
                <td><Activator handleActivate={() => this.props.handleActivate(this.props.id)} /></td>
                <td><Remover handleRemove={() => this.props.handleRemove(this.props.id)} /></td>
            </tr>
        );
    }
}


class ActivatedService extends Component {
    render() {
        return(
            <tr>
                <td><span>{ this.props.name }</span></td>
                <td><Deactivator handleStop={() => this.props.handleStop(this.props.name)} /></td>
            </tr>
        );
    }
}


export {Service, ActivatedService};