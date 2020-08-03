// An enum representing the three possible states of
// the neural engine running on the gateway
const Mode = Object.freeze({
    "Stopped": "Stopped",
    "Collecting": "Collecting",
    "Predicting": "Predicting"
})

const stateMapping = [
    Mode.Stopped,
    Mode.Collecting,
    Mode.Predicting
]


export default class NeuralState {
    constructor(mode, dataset) {
        this.mode = mode;
        this.dataset = dataset;
    }

    static FromAPI(apiResp) {
        var mode = stateMapping[apiResp.mode];
        var dataset = apiResp.dataset === "" ? "None" : apiResp.dataset;
        return new NeuralState(mode, dataset);
    }

    toReqFormat() {
        var mode = stateMapping.indexOf(this.mode);
        var dataset = this.dataset === "None" ? "" : this.dataset;
        return { mode: mode, dataset: dataset };
    }

    equals(other) {
        return this.mode === other.mode && this.dataset === other.dataset;
    }
}

export { Mode };