import React, { Component } from 'react';
import { ErrorMessage } from "./Error";
import { ActivatedService } from "./Service";
import "./List.css";


export default class ActivateList extends Component {
    constructor(props) {
        super(props);
        this.deactivateService = this.deactivateService.bind(this);
    }

    deactivateService(name) {
        this.props.handleServiceDeactivation(name);
    }

    render() {
        if (this.props.activatedServiceList.length === 0) {
            return (
                <div className="list">
                    <h2>Activated Services</h2>
                    <ErrorMessage name="No activated services!" />
                </div>
            );
        } else {
            return (
                <div className="list">
                    <h2>Activated Services</h2>
                    <table className="listTable">
                        <thead>
                            <tr>
                                <th>Name</th>
                                <th>Stop</th>
                            </tr>
                        </thead>
                        <tbody>
                            {this.props.activatedServiceList.map(activatedService => (
                                <ActivatedService key={activatedService}
                                    name={activatedService}
                                    handleStop={this.deactivateService} />
                            ))}
                        </tbody>
                    </table>
                </div>
            );
        }
    }
}

