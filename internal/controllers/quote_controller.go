package controllers

import "github.com/gofiber/fiber/v2"

func (s *Service) getAllQuoteRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupSearchRoute,
		s.setupRandomRoute,
		s.setupByCharacterRoute,
		s.setupByAudioIDRoute,
		s.setupCharactersRoute,
	}
}

func (s *Service) setupSearchRoute(routeGroup fiber.Router) {
	routeGroup.Get("/search", s.search)
}

func (s *Service) setupRandomRoute(routeGroup fiber.Router) {
	routeGroup.Get("/random", s.random)
}

func (s *Service) setupByCharacterRoute(routeGroup fiber.Router) {
	routeGroup.Get("/character/:id", s.byCharacter)
}

func (s *Service) setupByAudioIDRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/:audioId", s.byAudioID)
}

func (s *Service) setupCharactersRoute(routeGroup fiber.Router) {
	routeGroup.Get("/characters", s.characters)
}

func (s *Service) search(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'q' is required",
		})
	}

	lang := ctx.Query("lang", "en")
	limit := ctx.QueryInt("limit", 30)
	offset := ctx.QueryInt("offset", 0)
	characterID := ctx.Query("character")
	episode := ctx.QueryInt("episode", 0)
	forceFuzzy := ctx.QueryBool("fuzzy", false)

	response := s.QuoteService.Search(query, lang, limit, offset, characterID, episode, forceFuzzy)
	return ctx.JSON(fiber.Map{
		"query":   query,
		"results": response.Results,
		"total":   response.Total,
		"limit":   response.Limit,
		"offset":  response.Offset,
	})
}

func (s *Service) random(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "en")
	characterID := ctx.Query("character")
	episode := ctx.QueryInt("episode", 0)
	quote := s.QuoteService.Random(lang, characterID, episode)
	if quote == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no quotes available",
		})
	}
	return ctx.JSON(quote)
}

func (s *Service) byCharacter(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "en")
	characterID := ctx.Params("id")
	limit := ctx.QueryInt("limit", 50)
	offset := ctx.QueryInt("offset", 0)
	episode := ctx.QueryInt("episode", 0)

	response := s.QuoteService.GetByCharacter(lang, characterID, limit, offset, episode)
	return ctx.JSON(response)
}

func (s *Service) byAudioID(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "en")
	audioID := ctx.Params("audioId")

	quote := s.QuoteService.GetByAudioID(lang, audioID)
	if quote == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(quote)
}

func (s *Service) characters(ctx *fiber.Ctx) error {
	return ctx.JSON(s.QuoteService.GetCharacters())
}
