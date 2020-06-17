import React, {Component} from 'react';
import {ErrorMessage} from "./Error";
import { Service } from "./Service";
import "./List.css";


export default class ServiceList extends Component {
    constructor(props) {
        super(props);
        this.activateService = this.activateService.bind(this);
        this.removeService = this.removeService.bind(this);
    }

    activateService(name) {
        this.props.handleServiceActivation(name);
    }

    removeService(id) {
        this.props.handleServiceRemoval(id);
    }

    render() {
        if(this.props.serviceList.length === 0) {
            return(
                <div className="list">
                    <h2>Services</h2>
                    <ErrorMessage name="No services!" />
                </div>
            );
        }else{
            return(
                <div className="list">
                    <h2>Services</h2>
                    <table className="listTable">
                        <thead>
                        <tr>
                            <th>Name</th>
                            <th>Activate</th>
                            <th>Remove</th>
                        </tr>
                        </thead>
                        <tbody>
                        { this.props.serviceList.map(service => (
                            <Service key={service.id}
                                     id={service.id}
                                     name={service.name}
                                     handleActivate={this.activateService}
                                     handleRemove={this.removeService}/>
                        ))}
                        </tbody>
                    </table>
                </div>
            );
        }
    }
}
