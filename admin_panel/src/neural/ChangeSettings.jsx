import React, { Component } from "react";

export default class ChangeSettings extends Component {
    constructor(props) {
        super(props);

        this.state = {
            selectedMode: this.props.neuralState.mode,
            selectedDataset: this.props.neuralState.dataset,
        };

        this.modeChanged = this.modeChanged.bind(this);
        this.datasetChanged = this.datasetChanged.bind(this);
        this.forwardStateChange = this.forwardStateChange.bind(this);
    }

    modeChanged(evt) {
        this.setState({ selectedMode: evt.target.value })
    }

    datasetChanged(evt) {
        this.setState({ selectedDataset: evt.target.value })
    }

    forwardStateChange() {
        const { selectedMode, selectedDataset } = this.state;
        this.props.handleStateChange(selectedMode, selectedDataset);
    }


    render() {
        return (
            <div className="neuralSettings">
                <h2>Settings</h2>
                <div className="modeSetting">
                    <input type="radio" id="Stopped" name="mode" value="Stopped" onChange={(evt) => this.modeChanged(evt)} />
                    <label for="Stopped">Stopped</label>
                    <input type="radio" id="Collecting" name="mode" value="Collecting" onChange={(evt) => this.modeChanged(evt)} />
                    <label for="Collecting">Collecting</label>
                    <input type="radio" id="Predicting" name="mode" value="Predicting" onChange={(evt) => this.modeChanged(evt)} />
                    <label for="Predicting">Predicting</label>
                </div>

                <div className="datasetSetting">
                    <select name="datasetSelect" className="datasetSelect" onChange={(evt) => this.datasetChanged(evt)}>
                        <option value="" selected disabled hidden>Select a dataset to use</option>
                        {
                            this.props.datasetList.map(
                                dataset => <option key={dataset.name} value={dataset.name}>{dataset.name}</option>
                            )
                        }
                    </select>
                    <button onClick={this.forwardStateChange}>Set</button>
                </div>
            </div>
        );
    }
}