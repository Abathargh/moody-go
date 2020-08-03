import React, { Component } from "react";


export default class DatasetList extends Component {
    showKeys(keys) {
        alert(keys.toString());
    }

    render() {
        return (
            <div className="neuralBox">
                <h2>Datasets</h2>
                <div className="neuralTableBox">
                    <table className="neuralTable">
                        <thead>
                            <tr>
                                <th>Dataset</th>
                                <th>Keys</th>
                                <th>Remove</th>
                            </tr>
                        </thead>
                        <tbody>
                            {
                                this.props.datasetList.map(
                                    dataset =>
                                        <tr>
                                            <td>{dataset.name}</td>
                                            <td><button onClick={() => this.showKeys(dataset.keys)} >Show keys</button></td>
                                            <td><button className="remove" onClick={() => this.props.handleRemove(dataset.name)}>Remove</button></td>
                                        </tr>
                                )
                            }
                        </tbody>
                    </table>
                </div>
            </div >
        );
    }
}