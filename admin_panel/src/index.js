import React, { Component } from 'react';
import ReactDOM from 'react-dom';
import { Sidebar } from './sidebar.js';
import { MainCard } from './maincard.js';
import { BrowserRouter as Router } from 'react-router-dom';
import "./index.css";

class App extends Component {
    render() {
        return (
            <Router>
                <Sidebar />
                <MainCard />
            </Router>
        );
    }
}

ReactDOM.render(
    <App />,
    document.getElementById("root"),
)