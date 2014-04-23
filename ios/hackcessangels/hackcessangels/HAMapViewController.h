//
//  ViewController.h
//  HackcessAngels
//
//  Created by RIEUX Alexandre on 15/01/2014.
//  Copyright (c) 2014 RIEUX Alexandre. All rights reserved.
//


#import <UIKit/UIKit.h>

// model
#import "HATileOverlay.h"

// services
#import "HARestRequests.h"
#import "HAUserService.h"
#import "HACentralManager.h"

@interface HAMapViewController : UIViewController <MKMapViewDelegate>


/* Objet de la classe HARestRequest (à voir) */
@property (nonatomic, strong) HAUserService *userService;
@property (nonatomic, strong) HARestRequests *restRequests;
@property (nonatomic, strong) HATileOverlay *overlay;
@property (nonatomic, strong) HACentralManager *bluetoothmanager;

@property (nonatomic, weak) IBOutlet MKMapView *map;


@end