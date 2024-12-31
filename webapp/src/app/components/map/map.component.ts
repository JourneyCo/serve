import {Component, ViewChild} from '@angular/core';
import {GoogleMap, MapAdvancedMarker, MapInfoWindow } from "@angular/google-maps";

@Component({
  selector: 'app-map',
  templateUrl: './map.component.html',
  styleUrl: './map.component.css',
  standalone: true,
  imports: [GoogleMap, MapAdvancedMarker, MapInfoWindow]
})
export class MapComponent {
  @ViewChild(MapInfoWindow)
  infoWindow!: MapInfoWindow;

  options: google.maps.MapOptions = {
    center: {lat: 39.491482, lng: -104.874878},
    zoom: 12,
    mapId: '18c474b41c1ac65a',
  };
  display!: google.maps.LatLngLiteral;

  moveMap(event: google.maps.MapMouseEvent) {
    // @ts-ignore
    this.center = event.latLng.toJSON();
  }

  move(event: google.maps.MapMouseEvent) {
    // @ts-ignore
    this.display = event.latLng.toJSON();
  }

  advancedMarkerOptions: google.maps.marker.AdvancedMarkerElementOptions = {gmpDraggable: false};
  advancedMarkerPositions: google.maps.LatLngLiteral[] = [{lat: 39.491482, lng: -104.874878}];

  addAdvancedMarker(event: google.maps.MapMouseEvent) {
    // @ts-ignore
    this.advancedMarkerPositions.push(event.latLng.toJSON());
  }

  openInfoWindow(marker: MapAdvancedMarker) {
    this.infoWindow.open(marker);
  }

}
