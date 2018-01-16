import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { ApiService } from '../../service';
import { Observable } from 'rxjs/Observable';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Subject } from 'rxjs/Subject';
import { Subscription } from 'rxjs/Subscription';
import { MatDialogRef } from '@angular/material';

import 'rxjs/add/operator/take';
import 'rxjs/add/observable/interval';

@Component({
  selector: 'app-search-service',
  templateUrl: 'search-service.component.html',
  styleUrls: ['./search-service.component.scss'],
  encapsulation: ViewEncapsulation.None
})

export class SearchServiceComponent implements OnInit {
  searchStr = '';
  nodeAddr = '';
  seqs = [];
  timeOut = 3;
  resultTask: Subscription = null;
  totalResults: Array<Search> = [];
  results: Array<Search> = [];
  status = 0;
  private result: Subject<Array<Search>> = new BehaviorSubject<Array<Search>>([]);
  constructor(private api: ApiService, private dialogRef: MatDialogRef<SearchServiceComponent>) { }
  ngOnInit() {
    this.handle();
    // for (let index = 1; index <= 10; index++) {
    //   setTimeout(() => {
    //     this.result.next([{ seq: index, result: { node_key: `${index}`, apps: [`${index}-123456`] } }]);
    //   }, index * 1000);
    // }

    // setTimeout(() => {
    //   this.result.next([{ seq: 2, result: { node_key: 'b', apps: ['b123456'] } }]);
    // }, 1000);
    // setTimeout(() => {
    //   this.result.next([{ seq: 3, result: { node_key: 'a', apps: ['a123456'] } }]);
    // }, 2000);
    // setTimeout(() => {
    //   this.result.next([{ seq: 4, result: { node_key: 'b', apps: ['b123456'] } }]);
    // }, 3000);
    // setTimeout(() => {
    //   this.result.next([{ seq: 5, result: { node_key: 'a', apps: ['a123456', 'c654321'] } }]);
    // }, 4000);
    // setTimeout(() => {
    //   this.result.next([{ seq: 6, result: { node_key: 'c', apps: ['c123456'] } }]);
    // }, 5000);
    // setTimeout(() => {
    //   this.result.next([{ seq: 2, result: { node_key: 'b', apps: ['b123456', 'e123456'] } }]);
    // }, 6000);
    // setTimeout(() => {
    //   this.result.next([{ seq: 3, result: { node_key: 'a', apps: ['a123456', 'a654321'] } }]);
    // }, 7000);
    // setTimeout(() => {
    //   this.result.next([{ seq: 4, result: { node_key: 'b', apps: ['b123456', 'f123456', 'ff123456', 'xx123456'] } }]);
    // }, 8000);
    this.refresh();
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
      count: 1
    };
    data.append('client', 'socket');
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
    this.search();
    this.getResult();
    setTimeout(() => {
      this.status = 1;
    }, (this.timeOut + 6) * 1000);
  }
  search() {
    const data = new FormData();
    data.append('key', this.searchStr);
    Observable.interval(1000).take(this.timeOut).subscribe(() => {
      this.api.searchServices(this.nodeAddr, data).subscribe(seq => {
        this.seqs = this.seqs.concat(seq);
      });
    });
  }

  getResult() {
    Observable.interval(1000).take(this.timeOut + 5).subscribe(() => {
      this.api.getServicesResult(this.nodeAddr).subscribe(result => {
        this.result.next(result);
      });
    });
  }
  handle() {
    this.result.subscribe((results: Array<Search>) => {
      const tmp = this.filterSeq(results);
      this.unique(tmp);
      this.sort();
      this.results = this.totalResults;
    });
  }
  sort() {
    this.totalResults.sort((v1, v2) => {
      if (v1.seq < v2.seq) {
        return -1;
      }
      if (v1.seq > v2.seq) {
        return 1;
      }
      return 0;
    });
  }
  filterSeq(results: Array<Search>) {
    const tmpResults: Array<Search> = [];
    if (!results) {
      return;
    }
    results.forEach(result => {
      if (this.seqs.indexOf(result.seq) > - 1) {
        tmpResults.push(result);
      }
    });
    return tmpResults;
  }
  unique(results: Array<Search>) {
    if (!results) {
      return;
    }
    if (this.totalResults.length === 0) {
      this.totalResults = this.totalResults.concat(results);
      return;
    }
    for (let index = 0; index < results.length; index++) {
      for (let key = 0; key < this.totalResults.length; key++) {
        if (results[index].result.node_key === this.totalResults[key].result.node_key
          && results[index].result.apps.length === this.totalResults[key].result.apps.length) {
          results[index].result.apps.forEach(app => {
            if (this.totalResults[key].result.apps.indexOf(app) === -1) {
              this.totalResults[key].result.apps.push(app);
            }
          });
          break;
        }
        if (results[index].result.node_key !== this.totalResults[key].result.node_key) {
          let isExist = 0;
          isExist = this.totalResults.findIndex((el) => {
            if (el.result.node_key === results[index].result.node_key) {
              return true;
            }
            if (el.result.node_key !== results[index].result.node_key) {
              return false;
            }
            return true;
          });
          if (isExist < 0) {
            this.totalResults.push(results[index]);
          } else {
            this.totalResults.forEach((el, subIndex) => {
              if (results[index].result.node_key === el.result.node_key) {
                results[index].result.apps.forEach(app => {
                  if (el.result.apps.indexOf(app) === -1) {
                    this.totalResults[subIndex].result.apps.push(app);
                  }
                });
              }
            });
          }
          break;
        }
      }
    }
  }
}

export interface Search {
  result?: SearchResult;
  seq?: number;
}

export interface SearchResult {
  node_key?: string;
  apps?: Array<string>;
}
