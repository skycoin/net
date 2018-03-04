import { Component, OnInit, ViewEncapsulation, OnDestroy } from '@angular/core';
import { ApiService, AlertService } from '../../service';
import { Observable } from 'rxjs/Observable';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { Subscription } from 'rxjs/Subscription';
import { MatDialogRef, MatDialog } from '@angular/material';

import 'rxjs/add/operator/take';
import 'rxjs/add/observable/interval';

@Component({
  selector: 'app-search-service',
  templateUrl: 'search-service.component.html',
  styleUrls: ['./search-service.component.scss'],
  encapsulation: ViewEncapsulation.None
})

export class SearchServiceComponent implements OnInit, OnDestroy {
  searchStr = '';
  nodeAddr = '';
  seqs = [];
  timeOut = 1;
  resultTask: Subscription = null;
  searchTask: Subscription = null;
  totalResults: Array<SearchResultApp> = [];
  results: Array<SearchResultApp>;
  status = 0;
  SocketClient = 'socksc';
  private result: Subject<Array<Search>> = new BehaviorSubject<Array<Search>>([]);
  constructor(
    private api: ApiService,
    private alert: AlertService,
    private dialog: MatDialog,
    private dialogRef: MatDialogRef<SearchServiceComponent>) { }
  ngOnInit() {
    this.handle();
    this.refresh();
  }
  ngOnDestroy() {
    if (this.searchTask) {
      this.searchTask.unsubscribe();
    }
    if (this.resultTask) {
      this.resultTask.unsubscribe();
    }
    if (this.result) {
      this.result.unsubscribe();
    }
  }
  connectSocket(nodeKey: string, appKey: string) {
    if (!nodeKey || !appKey) {
      return;
    }
    const data = new FormData();
    const jsonStr = {
      label: '',
      nodeKey: nodeKey,
      appKey: appKey,
      count: 1,
      auto_start: false,
    };
    this.dialog.closeAll();
    this.alert.timer('connecting...', 30000);
    data.append('client', this.SocketClient);
    data.append('data', JSON.stringify(jsonStr));
    this.api.saveClientConnection(data).subscribe(res => {
      data.delete('data');
      data.delete('client');
    });
    data.append('toNode', nodeKey);
    data.append('toApp', appKey);
    this.api.connectSocketClicent(this.nodeAddr, data).subscribe(result => {
      console.log('conect socket client');
      this.dialogRef.close(result);
    });
  }
  refresh(ev?: Event) {
    if (ev) {
      ev.stopImmediatePropagation();
      ev.stopPropagation();
      ev.preventDefault();
    }
    this.status = 0;
    setTimeout(() => {
      this.status = 1;
    }, 3000);
    this.search();
    // this.getResult();
  }
  search() {
    const data = new FormData();
    data.append('key', this.searchStr);
    this.searchTask = Observable.interval(100).take(this.timeOut).subscribe(() => {
      this.api.searchServices(this.nodeAddr, data).subscribe(seq => {
        this.seqs = this.seqs.concat(seq);
        this.getResult();
      });
    }, err => {
      this.status = 1;
    });
  }

  getResult() {
    this.resultTask = Observable.interval(500).take(this.timeOut + 3).subscribe(() => {
      this.api.getServicesResult(this.nodeAddr).subscribe(result => {
        this.result.next(result);
        this.status = 1;
      });
    }, err => {
      this.status = 1;
    });
  }
  handle() {
    this.result.subscribe((results: Array<Search>) => {
      const tmp = this.filterSeq(results);
      if (!tmp) {
        return;
      }
      this.unique(tmp);
      this.sortByKey();
      this.results = this.totalResults;
    });
  }
  sortByKey() {
    this.totalResults.sort(
      function (a, b) {
        return a.node_key.localeCompare(b.node_key);
      });
  }
  trackByKey(index, app) {
    return app ? app.key : undefined;
  }
  filterSeq(results: Array<Search>) {
    const tmpResults: Array<Search> = [];
    if (!results) {
      return;
    }
    results.forEach(result => {
      const seqIndex = this.seqs.indexOf(result.seq);
      if (seqIndex > -1) {
        tmpResults.push(result);
      }
    });
    return tmpResults;
  }
  uniqueStep(results: Array<SearchResultApp>) {
    const res: Array<SearchResultApp> = [];
    const len = results.length;
    for (let i = 0; i < len; i++) {
      const item = results[i];
      let j = 0;
      let jLen = 0;
      for (j = 0, jLen = res.length; j < jLen; j++) {
        if (res[j].app_key === item.app_key && res[j].node_key === item.node_key) {
          break;
        }
      }
      if (j === jLen) {
        res.push(item);
      }
    }
    return res;
  }
  unique(results: Array<Search>) {
    if (results.length === 0) {
      return;
    }
    const apps: Array<SearchResultApp> = [];
    results.forEach(r => {
      r.result.forEach(app => {
        apps.push(app);
      });
    });
    this.seqs = [];
    this.totalResults = this.uniqueStep(apps);
    return;
  }
}

export interface Search {
  result?: Array<SearchResultApp>;
  seq?: number;
}
export interface SearchResultApp {
  node_key?: string;
  app_key?: string;
  location?: string;
}
export interface SearchResult {
  node_key?: string;
  apps?: Array<string>;
}
