import 'leaflet/dist/leaflet.css';
import './index.css';
import L from 'leaflet';
import {TileGrid} from '@liaozhai/tiles';

const pmax = 518788;
const amax = 2675.14;

function createStyle(feature) {
    const {population, area, cloudcover} = feature;
    const b = Math.ceil(255 * cloudcover / 100);
    return {fillStyle: `rgba(${[0,0,b].join(',')}, 0.8)`, strokeStyle: 'navy'};
};

window.addEventListener('load', () => {
    let map = L.map('map').setView([55.743733, 37.636413], 0);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    // const layerId = '82e314d6-1681-43a4-bb4d-f17a335a6527'
    const layerId = 'de615639-e60a-4dfd-8bc0-8e835a3f2345'
    // const layerId = '7905d39e-4529-4c52-8178-71eeafd7ec8e'
    const tg = new TileGrid({layerId, style: createStyle});
    tg.addTo(map);
});