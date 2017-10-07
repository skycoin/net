import { Component, OnInit } from '@angular/core';
import { ApiService, ConnData } from './service';
import { DataSource } from '@angular/cdk/collections';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/observable/of';
import { MdSnackBar } from '@angular/material';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
  listitems = ['Dashboard'];
  option = { selected: true };
  displayedColumns = ['index', 'type', 'key'];
  dataSource: ExampleDataSource = null;
  refreshTask = null;
  constructor(private api: ApiService, private snackBar: MdSnackBar) {
  }
  ngOnInit() {
    this.refresh();
    this.refreshTask = setInterval(() => {
      this.refresh();
    }, 5000);
  }
  TrackByKey(index, item) {
    return item ? item.index : undefined;
  }
  refresh(ev?: Event) {
    if (ev) {
      ev.stopImmediatePropagation();
      ev.stopPropagation();
      ev.preventDefault();
    }
    this.dataSource = new ExampleDataSource(this.api);
    if (ev) {
      this.snackBar.open('Refresh Successful', 'Dismiss', {
        duration: 2000,
        verticalPosition: 'top'
      });
    }
  }
  close() {
    clearInterval(this.refreshTask);
  }
}

export class ExampleDataSource extends DataSource<any> {
  constructor(private api: ApiService) {
    super();
  }
  connect(): Observable<ConnData[]> {
    return this.api.getAllNode();
  }

  disconnect() { }
}
