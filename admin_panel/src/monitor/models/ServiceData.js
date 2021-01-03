export default class ServiceData {
    constructor(service, data) {
        this.service = service;
        this.data = data;
    }

    static isServiceData(jsonObj) {
        return jsonObj.hasOwnProperty("service") && jsonObj.hasOwnProperty(("data"));
    }
}