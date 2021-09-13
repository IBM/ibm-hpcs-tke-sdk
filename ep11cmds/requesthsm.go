/*
 * Licensed Materials - Property of IBM
 *
 * 5900-A59
 *
 * OCO Source Materials
 *
 * (C) Copyright IBM Corp. 2021 All Rights Reserved
 *
 * The source code for this program is not published or other-
 * wise divested of its trade secrets, irrespective of what has
 * been deposited with the U.S. Copyright Office.
 */

// CHANGE HISTORY
//
// Date          Initials        Description

package ep11cmds

import (
	"github.com/IBM/ibm-hpcs-tke-sdk/common"
)

/*----------------------------------------------------------------------------*/
/* Creates Recovery HSMs in the case of the user requesting auto-init without */
/* any Recovery HSMs set up.                                                  */
/*                                                                            */
/* Input:                                                                     */
/* PluginContext -- contains the IAM access token and parameters identifying  */
/*    what resource group the user is working with                            */
/*                                                                            */
/* Output:                                                                    */
/* bool -- true for successful completion, false for error.  An error message */
/*    is always displayed before false is returned.                           */
/*----------------------------------------------------------------------------*/
func RequestHSM(authToken string, urlStart string, crypto_instance_id string, number_of_hsms int, hsm_type string) bool {

	req := common.CreateGetHsmsRequest(
		authToken, urlStart, crypto_instance_id)

	hsms, _, _, _, backup_region, err := common.SubmitQueryDomainsRequest(req)
	if err != nil {
		return false
	}
	if hsms == nil {
		return false
	}
	if len(hsms) == 0 {
		return false
	}

	if backup_region == "" {
		// Check available backup regions
		req, err := common.CreateQueryBackupRegionsHttpRequest(authToken, urlStart, hsm_type)
		if err != nil {
			return false
		}
		regions, err := common.SubmitAvailableBackupRegionsRequest(req)
		if err != nil {
			return false
		}
		if len(regions) == 0 {
			return false
		}
		backup_region = regions[0]
	}

	req2, err := common.CreateAssignHsmsHttpRequest(
		authToken, urlStart, crypto_instance_id, backup_region, number_of_hsms, hsm_type)
	if err != nil {
		return false
	}
	_, err = common.SubmitAssignHsmsRequest(req2)
	if err != nil {
		return false
	}

	return true
}
