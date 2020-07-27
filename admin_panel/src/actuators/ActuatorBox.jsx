import React, { Component } from "react";
import Mapping from "./Mapping"
import AddDeleteMapping from "./AddDeleteMapping";

export default class ActuatorBox extends Component {
    constructor(props) {
        super(props);

        this.forwardAdd = this.forwardAdd.bind(this);
        this.forwardDelete = this.forwardDelete.bind(this);
    }

    forwardAdd(situation, action) {
        this.props.handleAdd(this.props.ip, situation, action)
    }

    forwardDelete() {
        this.props.handleDelete(this.props.ip)
    }

    render() {
        return (
            <div className="actuatorBox">
                <h3>{this.props.ip}</h3>
                <AddDeleteMapping handleAdd={this.forwardAdd} handleDelete={this.forwardDelete} situationList={this.props.situationList} />
                <div className="actuatorBoxTable">
                    <table className="actTable">
                        <thead>
                            <tr>
                                <th>Situation</th>
                                <th>Action</th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                this.props.mappingList.map(
                                    mapping => (
                                        <Mapping key={mapping.situation}
                                            situation={mapping.situation}
                                            action={mapping.action} />
                                    )
                                )
                            }
                        </tbody>
                    </table>
                </div>
            </div>
        );
    }
}