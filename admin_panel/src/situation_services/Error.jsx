import React, { Component } from 'react';

// NIce centered div for an error in the format:
// <h1>Error type</h1>
// content ecc...

class Error extends Component {
    render() {
        return <div className="error"><h2>{this.props.name}</h2></div>;
    }
}

class ErrorMessage extends Component {
    render() {
        return <div className="ErrorMsg"><h2>{this.props.name}</h2></div>;
    }
}



export default Error;
export { ErrorMessage };