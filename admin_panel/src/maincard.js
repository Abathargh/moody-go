import React, { Component } from 'react'
import { Route, Switch } from 'react-router-dom';

import ServiceSituations from './situation_services'
import Monitor from "./monitor";
import Neural from "./neural"
import Actuators from "./actuators";

import "./maincard.css"

class MainCard extends Component {
    render() {
        return (
            <div className="maincard">
                <Switch>
                    <Route path="/services_situations">
                        <ServiceSituations />
                    </Route>
                    <Route path="/neural">
                        <Neural />
                    </Route>
                    <Route path="/actuators">
                        <Actuators />
                    </Route>
                    <Route path="/">
                        <Monitor />
                    </Route>
                </Switch>
            </div>
        );
    }
}

export { MainCard };