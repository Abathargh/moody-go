import React, {Component} from "react";
import ServiceBox from "./ServiceBox";

export default class ServiceGrid extends Component {
    render() {
        return(
            <div className="serviceGrid">
                {this.props.serviceList.map(service => <ServiceBox key={service.service} service={service.service} data={service.data}/>)}
            </div>
        );
    }
}