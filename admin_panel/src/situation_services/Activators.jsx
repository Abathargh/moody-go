import React, {Component} from "react";
import "./Activators.css"

class Activator extends Component {
    render() {
        return <button className="activate" onClick={this.props.handleActivate}>✓</button>;
    }
}

class Deactivator extends Component {
    render() {
        return <button className="deactivate" onClick={this.props.handleStop}>OFF</button>;
    }
}

class Creator extends Component {
    render() {
        return <button className="create" onClick={this.props.handleCreate}>New</button>;
    }
}

class Remover extends Component {
    render() {
        return <button className="remove" onClick={this.props.handleRemove}>X</button>;
    }
}

export {Activator, Deactivator, Creator, Remover};