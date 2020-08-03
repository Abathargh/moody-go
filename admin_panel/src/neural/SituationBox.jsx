import React, { Component } from "react";

export default class SituationBox extends Component {
    constructor(props) {
        super(props);

        this.state = {
            selectedSituation: this.props.currentSituation.id,
        }

        this.situationChanged = this.situationChanged.bind(this);
        this.forwardSituationChange = this.forwardSituationChange.bind(this);
    }

    situationChanged(evt) {
        this.setState({ selectedSituation: evt.target.value })
    }

    forwardSituationChange() {
        const situationChanged = this.state.selectedSituation;
        this.props.changeSituation(situationChanged);
    }

    render() {
        return (
            <div className="neuralBox situationBox">
                <div className="situationMonitor">
                    <h2>Situaton Monitor</h2>
                    <div className="stateTableBox">
                        <table className="stateTable">
                            <thead>
                                <tr>
                                    <th>Current Situation</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr>
                                    <td>{this.props.currentSituation.name}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
                <div className="situationSetting">
                    <h2>Settings</h2>
                    <p className="situationMsg"><i>The collected data tuples will be mapped with this situation</i></p>
                    <select name="situationSelect" className="situationSelect" onChange={(evt) => this.situationChanged(evt)}>
                        <option value="" selected disabled hidden>Select a situation to use</option>
                        {
                            this.props.situationList.map(
                                situation => <option key={situation.id} value={situation.id} >{situation.name}</option>
                            )
                        }
                    </select>
                    <button onClick={this.forwardSituationChange}>Set</button>
                </div>
            </div>
        );
    }
}