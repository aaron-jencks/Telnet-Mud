package crud

import (
	"fmt"
	"mud/actions/defined"
	"mud/actions/definitions"
	"mud/parsing_services/parsing"
	"mud/parsing_services/parsing/utils"
	"mud/parsing_services/player"
	crudUtils "mud/utils/crud"
	"mud/utils/handlers/crud"
	"net"
)

func createCrudAction(conn net.Conn, args []string,
	name, crudMethod string, validator CrudValidator,
	queryExecutor CrudExecutor, responseHandler CrudResponseHandler,
	reqModes []string) definitions.Action {
	return definitions.Action{
		Name:       fmt.Sprintf("%s %s", name, crudMethod),
		ValidModes: reqModes,
		Handler: func() parsing.CommandResponse {
			result := utils.GetDefaultInfoCommandResponse(conn)

			if !validator(args) {
				return result
			}

			nv := queryExecutor()
			responseHandler(nv)

			return result
		},
	}
}

// Creates an Action that wraps a Crud struct and calls it's Create method
func CreateCreateAction(conn net.Conn, args []string,
	name, usageString string, minArgs int, validator ArgumentValidator,
	argFmt ArgumentFormatter, respFmt ResponseFormatter,
	miscUpdate AfterActionHandler,
	reqModes []string, crudObj crudUtils.Crud) definitions.Action {
	username := player.GetConnUsername(conn)

	return createCrudAction(conn, args, name, "create",
		func(s []string) bool {
			return !(crud.CheckMinArgs(conn, args, minArgs, usageString) && validator(conn, args))
		}, func() interface{} {
			return crudObj.Create(argFmt(conn, args[1:])...)
		}, func(i interface{}) {
			player.EnqueueAction(username, defined.CreateInfoAction(conn, respFmt(i)))
		}, reqModes)
}

// Creates an Action that wraps a Crud struct and calls it's Retrieve method
func CreateRetrieveAction(conn net.Conn, args []string,
	name, usageString string, minArgs int, validator ArgumentValidator,
	argFmt RetrieveArgumentFormatter, respFmt ResponseFormatter,
	miscUpdate AfterActionHandler,
	reqModes []string, crudObj crudUtils.Crud) definitions.Action {
	username := player.GetConnUsername(conn)

	return createCrudAction(conn, args, name, "retrieve",
		func(s []string) bool {
			return !(crud.CheckMinArgs(conn, args, minArgs, usageString) && validator(conn, args))
		}, func() interface{} {
			return crudObj.Retrieve(argFmt(conn, args[1:]))
		}, func(i interface{}) {
			player.EnqueueAction(username, defined.CreateInfoAction(conn, respFmt(i)))
		}, reqModes)
}

// Creates an Action that wraps a Crud struct and calls it's Retrieve method
// Works with items that have multiple primary keys
func CreateMultiRetrieveAction(conn net.Conn, args []string,
	name, usageString string, minArgs int, validator ArgumentValidator,
	argFmt ArgumentFormatter, executor MultiKeyExecutor,
	respFmt ResponseFormatter,
	miscUpdate AfterActionHandler,
	reqModes []string, crudObj crudUtils.Crud) definitions.Action {
	username := player.GetConnUsername(conn)

	return createCrudAction(conn, args, name, "retrieve",
		func(s []string) bool {
			return !(crud.CheckMinArgs(conn, args, minArgs, usageString) && validator(conn, args))
		}, func() interface{} {
			return executor(conn, argFmt(conn, args[1:]))
		}, func(i interface{}) {
			player.EnqueueAction(username, defined.CreateInfoAction(conn, respFmt(i)))
		}, reqModes)
}

// Creates an Action that wraps a Crud struct and calls it's Update method
func CreateUpdateAction(conn net.Conn, args []string,
	name, usageString string, minArgs int,
	propertyIndex int, validator ArgumentValidator,
	argFmt RetrieveArgumentFormatter,
	valueUpdater UpdateNewValueFormatter,
	validPropertyNames []string,
	respFmt ResponseFormatter,
	miscUpdate AfterActionHandler,
	reqModes []string, crudObj crudUtils.Crud) definitions.Action {
	username := player.GetConnUsername(conn)

	return createCrudAction(conn, args, name, "update",
		func(s []string) bool {
			return !(crud.CheckMinArgs(conn, args, minArgs, usageString) && validator(conn, args) &&
				crud.CheckStringOptions(conn, args[propertyIndex], validPropertyNames, usageString, "property"))
		}, func() interface{} {
			ov := crudObj.Retrieve(argFmt(conn, args[1:]))
			nv := valueUpdater(ov, args[propertyIndex], args[propertyIndex+1:])
			return crudObj.Update(argFmt(conn, args[1:]), nv)
		}, func(i interface{}) {
			player.EnqueueAction(username, defined.CreateInfoAction(conn, respFmt(i)))
		}, reqModes)
}

// Creates an Action that wraps a Crud struct and calls it's Delete method
func CreateDeleteAction(conn net.Conn, args []string,
	name, usageString string, minArgs int, validator ArgumentValidator,
	argFmt RetrieveArgumentFormatter, respFmt ResponseFormatter,
	miscUpdate AfterActionHandler,
	reqModes []string, crudObj crudUtils.Crud) definitions.Action {
	username := player.GetConnUsername(conn)

	return createCrudAction(conn, args, name, "retrieve",
		func(s []string) bool {
			return !(crud.CheckMinArgs(conn, args, minArgs, usageString) && validator(conn, args))
		}, func() interface{} {
			ov := crudObj.Retrieve(argFmt(conn, args[1:]))
			crudObj.Delete(argFmt(conn, args[1:]))
			return ov
		}, func(i interface{}) {
			player.EnqueueAction(username, defined.CreateInfoAction(conn, respFmt(i)))
		}, reqModes)
}

// Creates an Action that wraps a Crud struct and calls it's Delete method
// Works with items that have multiple primary keys
func CreateMultiKeyDeleteAction(conn net.Conn, args []string,
	name, usageString string, minArgs int, validator ArgumentValidator,
	retriever MultiKeyExecutor, argFmt ArgumentFormatter, respFmt ResponseFormatter,
	miscUpdate AfterActionHandler,
	reqModes []string, crudObj crudUtils.Crud) definitions.Action {
	username := player.GetConnUsername(conn)

	return createCrudAction(conn, args, name, "retrieve",
		func(s []string) bool {
			return !(crud.CheckMinArgs(conn, args, minArgs, usageString) && validator(conn, args))
		}, func() interface{} {
			ov := retriever(conn, argFmt(conn, args[1:]))
			crudObj.Delete(argFmt(conn, args[1:])...)
			return ov
		}, func(i interface{}) {
			player.EnqueueAction(username, defined.CreateInfoAction(conn, respFmt(i)))
		}, reqModes)
}

func CreateCrudParser(name,
	createUsageString, retrieveUsageString, updateUsageString, deleteUsageString string,
	createMinArgs, retrieveMinArgs, updateMinArgs, deleteMinArgs int,
	createValidator, retrieveValidator, updateValidator, deleteValidator ArgumentValidator,
	createArgFmt ArgumentFormatter, retrievingFormatter RetrieveArgumentFormatter,
	createRespFmt, retrieveRespFmt, updateRespFmt, deleteRespFmt ResponseFormatter,
	createMiscUpdate, retrieveMiscUpdate, updateMiscUpdate, deleteMiscUpdate AfterActionHandler,
	validPropertyNames []string, propertyIndex int, valueUpdater UpdateNewValueFormatter,
	reqModes []string, crudObj crudUtils.Crud) parsing.CommandHandler {
	return func(conn net.Conn, args []string) {
		if crud.CrudChecks(conn, name, args) {
			return
		}

		username := player.GetConnUsername(conn)

		switch args[0] {
		case "create":
			player.EnqueueAction(username, CreateCreateAction(conn, args, name,
				createUsageString, createMinArgs, createValidator,
				createArgFmt, createRespFmt,
				createMiscUpdate,
				reqModes, crudObj,
			))

		case "retrieve":
			player.EnqueueAction(username, CreateRetrieveAction(conn, args, name,
				retrieveUsageString, retrieveMinArgs, retrieveValidator,
				retrievingFormatter, retrieveRespFmt,
				retrieveMiscUpdate,
				reqModes, crudObj,
			))

		case "update":
			player.EnqueueAction(username, CreateUpdateAction(conn, args, name,
				updateUsageString, updateMinArgs, propertyIndex, updateValidator,
				retrievingFormatter, valueUpdater,
				validPropertyNames,
				updateRespFmt,
				updateMiscUpdate,
				reqModes, crudObj,
			))

		case "delete":
			player.EnqueueAction(username, CreateDeleteAction(conn, args, name,
				deleteUsageString, deleteMinArgs, deleteValidator,
				retrievingFormatter, deleteRespFmt,
				deleteMiscUpdate,
				reqModes, crudObj,
			))
		}
	}
}

func CreateCrudParserMultiRetrieve(name,
	createUsageString, retrieveUsageString, deleteUsageString string,
	createMinArgs, retrieveMinArgs, deleteMinArgs int,
	createValidator, retrieveValidator, deleteValidator ArgumentValidator,
	createArgFmt ArgumentFormatter, deletingFormatter ArgumentFormatter,
	retriever MultiKeyExecutor,
	createRespFmt, retrieveRespFmt, deleteRespFmt ResponseFormatter,
	createMiscUpdate, retrieveMiscUpdate, deleteMiscUpdate AfterActionHandler,
	reqModes []string, crudObj crudUtils.Crud) parsing.CommandHandler {
	return func(conn net.Conn, args []string) {
		if crud.CrudChecks(conn, name, args) {
			return
		}

		username := player.GetConnUsername(conn)

		switch args[0] {
		case "create":
			player.EnqueueAction(username, CreateCreateAction(conn, args, name,
				createUsageString, createMinArgs, createValidator,
				createArgFmt, createRespFmt,
				createMiscUpdate,
				reqModes, crudObj,
			))

		case "retrieve":
			player.EnqueueAction(username, CreateMultiRetrieveAction(conn, args, name,
				retrieveUsageString, retrieveMinArgs, retrieveValidator,
				deletingFormatter, retriever, retrieveRespFmt,
				retrieveMiscUpdate,
				reqModes, crudObj,
			))

		case "update":
			player.EnqueueAction(username, defined.CreateInfoAction(conn, "Multi keyed object don't currently support updating"))

		case "delete":
			player.EnqueueAction(username, CreateMultiKeyDeleteAction(conn, args, name,
				deleteUsageString, deleteMinArgs, deleteValidator,
				retriever, deletingFormatter, deleteRespFmt,
				deleteMiscUpdate,
				reqModes, crudObj,
			))
		}
	}
}
