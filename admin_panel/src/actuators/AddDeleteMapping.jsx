import React, { Component } from "react";
import { AddMappingButton, DeleteMappingButton } from "./ButtonActions";
import { SituationListBox, NewActionBox } from "./InputSelect"

export default class AddDeleteMapping extends Component {
    constructor(props) {
        super(props);
        this.state = {
            selectedSituation: "",
            selectedAction: "",
        }
        this.setAction = this.setAction.bind(this);
        this.setSituation = this.setSituation.bind(this);
        this.forwardAdd = this.forwardAdd.bind(this);
        this.forwardDelete = this.forwardDelete.bind(this);
    }

    setAction(evt) {
        this.setState({ selectedAction: evt.target.value })
    }

    setSituation(evt) {
        this.setState({ selectedSituation: evt.target.value })
    }

    forwardAdd() {
        const { selectedSituation, selectedAction } = this.state;
        this.props.handleAdd(selectedSituation, selectedAction);
    }

    forwardDelete() {
        this.props.handleDelete();
    }

    render() {
        return (
            <div className="manageMapping">
                <SituationListBox situationList={this.props.situationList} handleChange={(evt) => this.setSituation(evt)} />
                <NewActionBox handleChange={(evt) => this.setAction(evt)} />
                <AddMappingButton handleAdd={this.forwardAdd} />
                <DeleteMappingButton handleDelete={this.forwardDelete} />
            </div>
        );
    }
}