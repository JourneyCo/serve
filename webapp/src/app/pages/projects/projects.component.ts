import {AfterViewInit, Component, inject, ViewChild} from '@angular/core';
import {MatPaginator, MatPaginatorModule} from '@angular/material/paginator';
import {MatSort, MatSortModule} from '@angular/material/sort';
import {MatTableDataSource, MatTableModule} from '@angular/material/table';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {Project} from "@models";
import {APIService} from "@services";
import {MapComponent, RegisterDialogComponent} from "@components";
import {DatePipe} from "@angular/common";
import {MatDialog} from "@angular/material/dialog";


@Component({
  selector: 'projects',
  styleUrl: 'projects.component.css',
  templateUrl: 'projects.component.html',
  imports: [MatFormFieldModule, MatInputModule, MatTableModule, MatSortModule, MatPaginatorModule, MapComponent,
  DatePipe],
})
export class ProjectsComponent implements AfterViewInit {
  displayedColumns: string[] = ['id', 'name', 'required', 'needed', 'date', 'created_at', 'updated_at', 'register']
  dataSource: MatTableDataSource<Project>;
  projects: Project[];
  clickedRow: Project | null;
  dialog = inject(MatDialog);
  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;

  constructor(
    private APIService: APIService,
  ) {
    this.APIService.getProjects().subscribe(data => {
      this.projects = data;
      this.dataSource = new MatTableDataSource(this.projects);
    });
  }

  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.sort;
  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  rowClick(row: Project) {
    if (this.clickedRow == row) {
      this.clickedRow = null
      return
    }
    this.clickedRow = row;
  }

  register(evt: any, row: Project) {
    console.log("clocked");
    const buttonElement = document.activeElement as HTMLElement; // Get the currently focused element
    buttonElement.blur();

    const dialogRef = this.dialog.open(RegisterDialogComponent, {
      data: row
    });

    dialogRef.afterClosed().subscribe(result => {
      console.log('The dialog was closed');
      // if (result !== undefined) {
      //   this.animal.set(result);
      // }
    });
    evt.stopPropagation();


  }
}
