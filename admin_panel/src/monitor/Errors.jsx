import React, {Component} from "react";
import "./boxes.css"

class Error extends Component {
    render() {
        return <div className="error"><h2>Can't fetch service data!</h2></div>;
    }
}

class Loading extends Component {
    render() {
        return <div className="error"><h2>Loading...</h2></div>;
    }
}

class Empty extends Component {
    render() {
        return <div className="error"><h2>No active service!</h2></div>;
    }
}

export {Empty, Error, Loading};