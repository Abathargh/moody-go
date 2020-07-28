import React, { Component } from "react";


export default class DatasetList extends Component {
    render() {
        return (
            <div className="datasetList">
                <ul>
                    {
                        this.props.datasetList.map(
                            dataset => <li>{dataset.name}: {dataset.keys.toString()}</li>
                        )
                    }
                </ul>
            </div>
        );
    }
}