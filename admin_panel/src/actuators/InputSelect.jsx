import React, { Component } from "react"


class SituationListBox extends Component {
    render() {
        return (
            <div className="mappingOption">
                <select name="situationList" className="inputBox" onChange={this.props.handleChange}>
                    <option value="" selected disabled hidden>Select a Situation</option>
                    {
                        this.props.situationList.map(
                            situation => <option key={situation.id} value={situation.id}>{situation.name}</option>
                        )
                    }
                </select>
            </div>
        )
    }
}

class NewActionBox extends Component {
    render() {
        return <div className="mappingOption"><input className="inputBox" type="text" placeholder="Action" onChange={this.props.handleChange} /></div>;
    }
}



export { SituationListBox, NewActionBox };