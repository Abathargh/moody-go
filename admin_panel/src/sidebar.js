import React, { Component } from 'react'
import { Link } from 'react-router-dom';

import './sidebar.css'
import logo from './logo.png'

class Sidebar extends Component {
    render() { 
        return(
            <div className="sidebar">
                <img src={logo} className="logo" alt="logo" />
                <ul>
                    <li><Link to="/">Monitor</Link></li>
                    <li><Link to="/services_situations">Services &amp; Situations</Link></li>
                    <li><Link to="/neural">Neural</Link></li>
                    <li><Link to="/actuators">Manage Actuators</Link></li>
                </ul>
            </div>
        );
    }
}

export { Sidebar };