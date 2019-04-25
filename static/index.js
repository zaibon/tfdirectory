Vue.component('status-meaning', {
    template: '#statusHelpModal-template'
})

Vue.component("node-detail", {
    template: "#node-detail-template",
    props: ['node'],
})

const baseURL = 'https://3aab148f.ngrok.io'
// const baseURL = 'http://localhost:8081'
var app = new Vue({
    el: '#app',
    data: {
        showStatusHelpModal: false,
        resource: {
            cru:0,
            mru:0,
            hru:0,
            sru:0,
        },
        nodes: [],
        farmers: [],
        selectedFarmer: "",
        selectedCountry: "",
        selectedNodeID: null
    },
    created: function(){
        this.fetchFarmers()
        this.fetchNodes()
    },
    filters: {
        nodeStatus: function(value){
            if (!value) {
                return 'down'
            }

            updated = Date.parse(value)
            now = Date.now()
            delta = now - updated

            if (delta <= 600000){ // 10 minutes or less
                return 'up'
            }
            if (600000 < delta && delta < 900000) { //between 10 and 15 minutes
                return 'likely down'
            }
            return 'down'
        },
        statusColor: function(value){
            if (!value) {
                return 'badge-danger'
            }

            updated = Date.parse(value)
            now = Date.now()
            delta = now - updated

            if (delta <= 600000){ // 10 minutes or less
                return 'badge-success'
            }
            if (600000 < delta && delta < 900000) { //between 10 and 15 minutes
                return 'badge-warning '
            }
            return 'badge-danger'
        }
    },
    computed: {
        countries: function() {
            var countries = new Set()
            this.nodes.forEach(node => {
                countries.add(node.location.country)
            })
            if (this.selectedFarmer == "") {
                this.farmers.forEach(farmer => {
                    countries.add(farmer.location.country)
                })
            }
            return Array.from(countries).sort();
        },
        filteredNodes: function() {
            var self = this
            return this.nodes.filter(function(node){
                if (self.selectedCountry != "" && node.location.country != self.selectedCountry) {
                    return false
                }
                if (self.selectedFarmer != "" && node.farmer_id != self.selectedFarmer) {
                    return false
                }
                if (self.resource.cru != 0 && node.total_resources.cru < self.resource.cru) {
                    return false
                }
                if (self.resource.mru != 0 && node.total_resources.mru < self.resource.mru) {
                    return false
                }
                if (self.resource.hru != 0 && node.total_resources.hru < self.resource.hru) {
                    return false
                }
                if (self.resource.sru != 0 && node.total_resources.sru < self.resource.sru) {
                    return false
                }
                return true
            })
        },
        selectedNode: function(){
            for (let i = 0; i < this.nodes.length; i++) {
                const node = this.nodes[i];
                if (node.node_id == this.selectedNodeID) {
                    this.farmers.forEach(farmer => {
                        if (farmer.iyo_organization == node.farmer_id){
                            node.farmer = farmer
                        }
                    });

                    node.available_resources = {
                        cru: node.used_resources.cru ? node.total_resources.cru - node.reserved_resources.cru : 0,
                        mru: node.used_resources.mru ? node.total_resources.mru - node.reserved_resources.mru: 0,
                        hru: node.used_resources.hru ? node.total_resources.hru - node.reserved_resources.hru: 0,
                        sru: node.used_resources.sru ? node.total_resources.sru - node.reserved_resources.sru: 0,
                    }
                   return node
                }
            }
            return null
        }
    },
    methods: {
        reset: function() {
            this.selectedCountry = ""
            this.selectedFarmer = ""
            this.resource = {
                cru:0,
                mru:0,
                hru:0,
                sru:0,
            }
            // this.nodes = []
        },
        fetchNodes: function(){
            var xhr = new XMLHttpRequest()
            var self = this
            var url = baseURL+"/api/nodes"
            // function encodeQueryData(data) {
            //     const ret = [];
            //     for (let d in data){
            //         if (data != "") {
            //             ret.push(encodeURIComponent(d) + '=' + encodeURIComponent(data[d]));
            //         }
            //     }
            //     return ret.join('&');
            //  }
            // url += "?" + encodeQueryData({
            //     "country": self.selectedCountry,
            //     "farmer": self.selectedFarmer,
            //     "cru": self.resource.cru,
            //     "mru": self.resource.mru,
            //     "hru": self.resource.hru,
            //     "sru": self.resource.sru
            // })
            xhr.open('GET', url)
            xhr.onload = function () {
                self.nodes = JSON.parse(xhr.responseText)
                self.nodes.forEach(node => {
                    if (node.location == null){
                        node.location = {
                            "continent": "unknown",
                            "country": "unknown",
                            "city": "unknown",
                        }
                    }
                });
            }
            xhr.send()
        },
        fetchFarmers: function(){
            var xhr = new XMLHttpRequest()
            var self = this
            var url = baseURL+"/api/farmers"
            xhr.open('GET', url)
            xhr.onload = function () {
                self.farmers = JSON.parse(xhr.responseText)
                self.farmers.sort((a,b) => (a.name > b.name) ? 1 : ((b.name > a.name) ? -1 : 0));
            }
            xhr.send()
        }
    },
});
