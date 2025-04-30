import {MapAdvancedMarker} from '@angular/google-maps';
import {Project} from './project.model';

export interface Marker {
    position: google.maps.LatLngLiteral;
    title: string;
    label: string;
}

export type ServeMarker = MapAdvancedMarker & {project: Project}
