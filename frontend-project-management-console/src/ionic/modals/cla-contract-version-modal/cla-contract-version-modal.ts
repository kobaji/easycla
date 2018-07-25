import { Component } from '@angular/core';
import { IonicPage,  ModalController, NavController, NavParams, ViewController, } from 'ionic-angular';
import { PopoverController } from 'ionic-angular';
import { ClaService } from 'cla-service';

@IonicPage({
  segment: 'cla-contract-version-modal'
})
@Component({
  selector: 'cla-contract-version-modal',
  templateUrl: 'cla-contract-version-modal.html',
})
export class ClaContractVersionModal {

  claProjectId: string;
  documentType: string; // individual | corporate
  documents: any;
  currentDocument: any;
  previousDocuments: any;

  constructor(
    public navCtrl: NavController,
    public navParams: NavParams,
    public viewCtrl: ViewController,
    private popoverCtrl: PopoverController,
    public modalCtrl: ModalController,
    private claService: ClaService,
  ) {
    this.claProjectId = this.navParams.get('claProjectId');
    this.documentType = this.navParams.get('documentType');
    this.documents = this.navParams.get('documents').reverse();
    if (this.documents.length > 0) {
      this.currentDocument = this.documents.slice(0, 1);
      console.log("currentDoc");
      console.log(this.currentDocument);
      if (this.documents.length > 1) {
        this.previousDocuments = this.documents.slice(1);
      }
    }

    console.log(this.documents);
    this.getDefaults();
  }

  getDefaults() {

  }

  ngOnInit() {

  }

  contractPopover(ev, document) {
    let currentContractActions = {
      items: [
        {
          label: 'View PDF',
          callback: 'contractView',
          callbackData: {
            document: document,
          }
        },
      ]
    };
    let popover = this.popoverCtrl.create(
      'ActionPopoverComponent',
      currentContractActions,
    );

    popover.present({
      ev: ev
    });

    popover.onDidDismiss((popoverData) => {
      if(popoverData) {
        this.popoverResponse(popoverData);
      }
    });
  }

  /**
   * Called if popover dismissed with data. Passes data to a callback function
   * @param  {object} popoverData should contain .callback and .callbackData
   */
  popoverResponse(popoverData) {
    let callback = popoverData.callback;
    if(this[callback]) {
      this[callback](popoverData.callbackData);
    }
  }

  contractView(data) {
    console.log('contract view');
    console.log(data);
    let majorVersion = data.document.document_major_version;
    let minorVersion = data.document.document_minor_version;
    let url = this.claService.getProjectDocumentRevisionPdf(this.claProjectId, this.documentType, majorVersion, minorVersion);
    console.log(url);
    window.open(url, '_blank');
  }

  openClaContractUploadModal() {
    let modal = this.modalCtrl.create('ClaContractUploadModal', {
      claProjectId: this.claProjectId,
      documentType: this.documentType,
    });
    modal.present();
  }

  dismiss() {
    this.viewCtrl.dismiss();
  }

}