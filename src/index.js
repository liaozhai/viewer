import 'leaflet/dist/leaflet.css';
import './index.css';
import L from 'leaflet';
import {TileGrid} from '@liaozhai/tiles';

const pmax = 518788;
const amax = 2675.14;

function createStyle(feature) {
    const {population, area, cloudcover} = feature;
    const a = area / amax;
    const p = population / pmax;
    return {fillStyle: `rgba(${[255,0,0].join(',')}, ${a})`, strokeStyle: 'navy'};
};

window.addEventListener('load', () => {
    const stream = new EventSource("/stream");
    stream.addEventListener("message", function(e){
        console.log(e.data);    
    });

    let map = L.map('map').setView([55.743733, 37.636413], 6);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    const layerId = '82e314d6-1681-43a4-bb4d-f17a335a6527'
    // const layerId = 'de615639-e60a-4dfd-8bc0-8e835a3f2345'
    // const layerId = '7905d39e-4529-4c52-8178-71eeafd7ec8e'
    const tg = new TileGrid({layerId, style: createStyle});
    tg.addTo(map);
});