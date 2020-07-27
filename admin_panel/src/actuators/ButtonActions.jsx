import React, { Component } from "react"

class AddMappingButton extends Component {
    render() {
        return <div className="mappingOption"><button onClick={this.props.handleAdd}>Add Mapping</button></div>;
    }
}

class DeleteMappingButton extends Component {
    render() {
        return <div className="mappingOption"><button onClick={this.props.handleDelete}>Delete Mappings</button></div>;
    }
}



export { AddMappingButton, DeleteMappingButton };