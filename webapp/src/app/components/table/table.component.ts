import { AsyncPipe, DecimalPipe } from '@angular/common';
import { Component, QueryList, ViewChildren } from '@angular/core';
import { Observable } from 'rxjs';
import { Project } from './project';
import { ProjectService } from '../../services/project.service';
import { NgbdSortableHeader, SortEvent } from '../../directives/sortable.directive';
import { FormsModule } from '@angular/forms';
import { NgbHighlight, NgbPaginationModule } from '@ng-bootstrap/ng-bootstrap';

@Component({
  selector: 'app-table',
  imports: [DecimalPipe, FormsModule, AsyncPipe, NgbHighlight, NgbdSortableHeader, NgbPaginationModule],
  templateUrl: './table.component.html',
  providers: [ProjectService, DecimalPipe],
})
export class TableComponent {
  projects$: Observable<Project[]>;
  total$: Observable<number>;

  @ViewChildren(NgbdSortableHeader) headers: QueryList<NgbdSortableHeader> | undefined;

  constructor(public service: ProjectService) {
    this.projects$ = service.projects$;
    this.total$ = service.total$;
  }

  onSort({ column, direction }: SortEvent) {
    // resetting other headers
    // @ts-ignore
    this.headers.forEach((header) => {
      if (header.sortable !== column) {
        header.direction = '';
      }
    });

    this.service.sortColumn = column;
    this.service.sortDirection = direction;
  }

  rowClicked(id: number) {
    console.log(id);
  }

  register(id: number) {
    console.log("registering! : ", id)
  }
}
