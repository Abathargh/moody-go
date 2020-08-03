import React, { Component } from "react"

export default class NewDatasetBox extends Component {
    constructor(props) {
        super(props);

        this.state = {
            newDatasetName: "",
        }

        this.changeHandler = this.changeHandler.bind(this);
        this.forwardSituationChange = this.forwardSituationChange.bind(this);
    }

    changeHandler(evt) {
        this.setState({ newDatasetName: evt.target.value })
    }

    forwardSituationChange() {
        const newDatasetName = this.state.newDatasetName;
        this.props.handleCreate(newDatasetName)
    }


    render() {
        return (
            <div className="neuralBox">
                <h2>Create Dataset</h2>
                <div className="divMsg">
                    <p><i>
                        When creating a new dataset, the currently active services will
                        be used as the dataset keys.
                        </i></p>

                    <p><i>
                        The dataset will then be used only if all the given services are active
                        and are receiving data.
                    </i></p>
                </div>
                <div className="createDatasetBox">
                    <input type="text" onChange={(evt) => this.changeHandler(evt)} placeholder="Insert the new dataset name"></input>
                    <button onClick={this.forwardSituationChange}>Set</button>
                </div>
            </div>
        );
    }
}