import React, {Component } from 'react';

// NIce centered div for an error in the format:
// <h1>Error type</h1>
// content ecc...

class Loading extends Component {
    render() {
        return <div className="Error"><h2>{this.props.name}</h2></div>;
    }
}

export default Loading;