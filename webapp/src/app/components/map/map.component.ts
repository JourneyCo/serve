import { Component } from '@angular/core';
import {GoogleMap, MapAdvancedMarker } from "@angular/google-maps";

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrl: './map.component.css',
  standalone: true,
  imports: [GoogleMap, MapAdvancedMarker]
})
export class MapComponent {

  options: google.maps.MapOptions = {
    center: {lat: 39.491482, lng: -104.874878},
    zoom: 12,
    mapId: '18c474b41c1ac65a',
  };
  display: google.maps.LatLngLiteral | undefined;

  moveMap(event: google.maps.MapMouseEvent) {
    // @ts-ignore
    this.center = (event.latLng.toJSON());
  }

  move(event: google.maps.MapMouseEvent) {
    // @ts-ignore
    this.display = event.latLng.toJSON();
  }

  advancedMarkerOptions: google.maps.marker.AdvancedMarkerElementOptions = {gmpDraggable: false};
  advancedMarkerPositions: google.maps.LatLngLiteral[] = [];

  addAdvancedMarker(event: google.maps.MapMouseEvent) {
    // @ts-ignore
    this.advancedMarkerPositions.push(event.latLng.toJSON());
  }
}
