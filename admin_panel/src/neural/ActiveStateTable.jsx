import React, { Component } from "react"

export default class ActiveStateTable extends Component {
    render() {
        return (
            <div className="stateTableBox">
                <table className="stateTable">
                    <thead>
                        <tr>
                            <th>State</th>
                            <th>Dataset in use</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>{this.props.neuralState.mode}</td>
                            <td>{this.props.neuralState.mode === "Stopped" ? "None" : this.props.neuralState.dataset}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        );
    }
}