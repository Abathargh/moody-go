export default class IpMappingReq {
    constructor(ip, situation, action) {
        this.ip = ip;
        this.situation = parseInt(situation);
        this.action = parseInt(action);
    }
}