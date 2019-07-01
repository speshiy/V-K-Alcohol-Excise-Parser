package cscript

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/_main/models/mscript"
	"github.com/speshiy/V-K-Alcohol-Excise-Parse/common"
)

//GetScriptTypes return all script types by companyID
func GetScriptTypes(c *gin.Context) {
	var err error
	var scriptTypes []mscript.ScriptType
	var companyID uint

	companyID, err = common.ParseParam(c, "CompanyID")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = mscript.GetScriptTypes(c, companyID, &scriptTypes)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "data": scriptTypes})
}

//PostScriptType add new script type
func PostScriptType(c *gin.Context) {
	var err error
	var scriptType mscript.ScriptType

	if err := c.ShouldBindJSON(&scriptType); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = scriptType.Post(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Done", "data": scriptType})
}

//PutScriptType update script type
func PutScriptType(c *gin.Context) {
	var err error
	var scriptType mscript.ScriptType

	if err := c.ShouldBindJSON(&scriptType); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	err = scriptType.Put(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Done", "data": scriptType})
}

//DeleteScriptType delete script type
func DeleteScriptType(c *gin.Context) {
	var err error
	var id uint
	var scriptType mscript.ScriptType

	id, err = common.ParseParam(c, "ID")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	if id < 1 {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": "Script type ID invalid"})
		return
	}

	scriptType.ID = id
	err = scriptType.Delete(c, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "false", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "true", "message": "Done"})
}
