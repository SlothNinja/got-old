package main

import (
	"net/http/httptest"

	"bitbucket.org/SlothNinja/user"
	"github.com/gin-gonic/gin"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("g.grid.lampAreas", func() {
	var (
		c                 *gin.Context
		resp              *httptest.ResponseRecorder
		g                 game
		u1, u2            user.User2
		expectedAreas, as []area
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()
		g = createGame(c, u1, u2)
	})

	Context("when no area selected", func() {
		It("should return no areas", func() {
			Expect(g.grid.lampAreas(noArea)).To(HaveLen(0))
		})
	})

	Context("when full grid without any thieves", func() {

		type entry struct {
			selected areaID
			thief    areaID
			expected []areaID
		}

		DescribeTable("when valid selected area",
			func(e entry) {
				if e.thief != noAreaID {
					thiefArea := g.grid.area(e.thief.row, e.thief.column)
					thiefArea.thief.pid = 1
					g.grid = g.grid.updateArea(thiefArea)
				}

				from := g.grid.area(e.selected.row, e.selected.column)
				expectedAreas = make([]area, len(e.expected))
				for i, aid := range e.expected {
					expectedAreas[i] = g.grid.area(aid.row, aid.column)
				}
				as = g.grid.lampAreas(from)
				Expect(as).To(ConsistOf(expectedAreas))
			},
			Entry("area(1,1)", entry{
				selected: areaID{1, 1},
				thief:    noAreaID,
				expected: []areaID{areaID{1, 8}, areaID{6, 1}},
			}),
			Entry("area(1,2), thief south", entry{
				selected: areaID{1, 2},
				thief:    areaID{2, 2},
				expected: []areaID{areaID{1, 1}, areaID{1, 8}},
			}),
			Entry("area(2,3), thief north", entry{
				selected: areaID{2, 3},
				thief:    areaID{1, 3},
				expected: []areaID{areaID{2, 1}, areaID{2, 8}, areaID{6, 3}},
			}),
			Entry("area(4,4), thief east", entry{
				selected: areaID{4, 4},
				thief:    areaID{4, 5},
				expected: []areaID{areaID{4, 1}, areaID{1, 4}, areaID{6, 4}},
			}),
			Entry("area(5,5), thief west", entry{
				selected: areaID{5, 5},
				thief:    areaID{5, 4},
				expected: []areaID{areaID{1, 5}, areaID{6, 5}, areaID{5, 8}},
			}),
		)
	})
})

var _ = Describe("g.swordAreasFor", func() {
	var (
		c                 *gin.Context
		cp                player
		resp              *httptest.ResponseRecorder
		g                 game
		cu, u1, u2        user.User2
		expectedAreas, as []area
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()
		g = createGame(c, u1, u2)

		if g.CPUserIndices[0] == 1 {
			cu = u1
		} else {
			cu = u2
		}
		user.WithCurrent(c, cu)

		cp = g.currentPlayerFor(cu)
	})

	Context("when no area selected", func() {
		It("should return no areas", func() {
			Expect(g.grid.swordAreasFor(cp, noArea)).To(HaveLen(0))
		})
	})

	Context("when full grid without any thieves", func() {
		var ta area

		type entry struct {
			selected areaID
			thief    areaID
			opponent bool
			expected []areaID
		}

		DescribeTable("when valid selected area",
			func(e entry) {
				from := g.grid.area(e.selected.row, e.selected.column)
				ta = g.grid.area(e.thief.row, e.thief.column)
				ta.thief.pid = cp.id
				if e.opponent {
					ta.thief.pid += 1
				}
				g.grid = g.grid.updateArea(ta)

				expectedAreas = make([]area, len(e.expected))
				for i, aid := range e.expected {
					expectedAreas[i] = g.grid.area(aid.row, aid.column)
				}
				as = g.grid.swordAreasFor(cp, from)
				Expect(as).To(ConsistOf(expectedAreas))
			},
			Entry("east", entry{
				selected: areaID{1, 1},
				thief:    areaID{1, 2},
				opponent: true,
				expected: []areaID{areaID{1, 2}},
			}),
			Entry("can't bump east", entry{
				selected: areaID{1, 1},
				thief:    areaID{1, 8},
				opponent: true,
				expected: noAreaIDS,
			}),
			Entry("west", entry{
				selected: areaID{1, 3},
				thief:    areaID{1, 2},
				opponent: true,
				expected: []areaID{areaID{1, 2}},
			}),
			Entry("can't bump west", entry{
				selected: areaID{1, 3},
				thief:    areaID{1, 1},
				opponent: true,
				expected: noAreaIDS,
			}),
			Entry("north", entry{
				selected: areaID{6, 2},
				thief:    areaID{2, 2},
				opponent: true,
				expected: []areaID{areaID{2, 2}},
			}),
			Entry("can't bump north", entry{
				selected: areaID{5, 3},
				thief:    areaID{1, 3},
				opponent: true,
				expected: noAreaIDS,
			}),
			Entry("south", entry{
				selected: areaID{2, 5},
				thief:    areaID{4, 5},
				opponent: true,
				expected: []areaID{areaID{4, 5}},
			}),
			Entry("can't bump south", entry{
				selected: areaID{5, 4},
				thief:    areaID{6, 4},
				opponent: true,
				expected: noAreaIDS,
			}),
		)
	})
})

var _ = Describe("g.carpetAreas", func() {
	var (
		c                 *gin.Context
		cp                player
		resp              *httptest.ResponseRecorder
		g                 game
		cu, u1, u2        user.User2
		expectedAreas, as []area
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()
		g = createGame(c, u1, u2)

		if g.CPUserIndices[0] == 1 {
			cu = u1
		} else {
			cu = u2
		}
		user.WithCurrent(c, cu)

		cp = g.currentPlayerFor(cu)
	})

	Context("when no area selected", func() {
		It("should return no areas", func() {
			Expect(g.grid.carpetAreas(noArea)).To(HaveLen(0))
		})
	})

	Context("when full grid without any thieves", func() {
		var ta area

		type entry struct {
			selected areaID
			thief    areaID
			opponent bool
			empty    []areaID
			expected []areaID
		}

		DescribeTable("when valid selected area",
			func(e entry) {
				from := g.grid.area(e.selected.row, e.selected.column)
				ta = g.grid.area(e.thief.row, e.thief.column)
				ta.thief.pid = cp.id
				if e.opponent {
					ta.thief.pid += 1
				}
				g.grid = g.grid.updateArea(ta)

				for _, aid := range e.empty {
					empty := noArea
					empty.areaID = aid
					g.grid = g.grid.updateArea(empty)
				}

				expectedAreas = make([]area, len(e.expected))
				for i, aid := range e.expected {
					expectedAreas[i] = g.grid.area(aid.row, aid.column)
				}
				as = g.grid.carpetAreas(from)
				Expect(as).To(ConsistOf(expectedAreas))
			},
			Entry("east thief", entry{
				selected: areaID{1, 1},
				thief:    areaID{1, 2},
				opponent: true,
				empty:    noAreaIDS,
				expected: noAreaIDS,
			}),
			Entry("east empty", entry{
				selected: areaID{1, 1},
				thief:    areaID{1, 8},
				opponent: true,
				empty:    []areaID{areaID{1, 2}},
				expected: []areaID{areaID{1, 3}},
			}),
			Entry("east empty", entry{
				selected: areaID{1, 1},
				thief:    areaID{1, 8},
				opponent: true,
				empty:    []areaID{areaID{1, 2}, areaID{1, 3}},
				expected: []areaID{areaID{1, 4}},
			}),
			Entry("east empty, but thief", entry{
				selected: areaID{1, 1},
				thief:    areaID{1, 5},
				opponent: true,
				empty:    []areaID{areaID{1, 2}, areaID{1, 3}, areaID{1, 4}},
				expected: noAreaIDS,
			}),
			Entry("west thief", entry{
				selected: areaID{2, 3},
				thief:    areaID{2, 2},
				opponent: true,
				empty:    noAreaIDS,
				expected: noAreaIDS,
			}),
			Entry("west empty", entry{
				selected: areaID{2, 8},
				thief:    areaID{2, 1},
				opponent: true,
				empty:    []areaID{areaID{2, 7}},
				expected: []areaID{areaID{2, 6}},
			}),
			Entry("west empty", entry{
				selected: areaID{2, 8},
				thief:    areaID{2, 1},
				opponent: true,
				empty:    []areaID{areaID{2, 7}, areaID{2, 6}},
				expected: []areaID{areaID{2, 5}},
			}),
			Entry("west empty, but thief", entry{
				selected: areaID{4, 5},
				thief:    areaID{4, 1},
				opponent: true,
				empty:    []areaID{areaID{4, 4}, areaID{4, 3}, areaID{4, 2}},
				expected: noAreaIDS,
			}),
			Entry("down thief", entry{
				selected: areaID{4, 4},
				thief:    areaID{5, 4},
				opponent: true,
				empty:    noAreaIDS,
				expected: noAreaIDS,
			}),
			Entry("down empty", entry{
				selected: areaID{4, 5},
				thief:    areaID{2, 1},
				opponent: true,
				empty:    []areaID{areaID{5, 5}},
				expected: []areaID{areaID{6, 5}},
			}),
			Entry("down empty", entry{
				selected: areaID{2, 6},
				thief:    areaID{2, 1},
				opponent: true,
				empty:    []areaID{areaID{3, 6}, areaID{4, 6}},
				expected: []areaID{areaID{5, 6}},
			}),
			Entry("down empty, but thief", entry{
				selected: areaID{2, 5},
				thief:    areaID{6, 5},
				opponent: true,
				empty:    []areaID{areaID{3, 5}, areaID{4, 5}, areaID{5, 5}},
				expected: noAreaIDS,
			}),
			Entry("up thief", entry{
				selected: areaID{6, 4},
				thief:    areaID{5, 4},
				opponent: true,
				empty:    noAreaIDS,
				expected: noAreaIDS,
			}),
			Entry("up empty", entry{
				selected: areaID{6, 5},
				thief:    areaID{2, 1},
				opponent: true,
				empty:    []areaID{areaID{5, 5}},
				expected: []areaID{areaID{4, 5}},
			}),
			Entry("up empty", entry{
				selected: areaID{5, 6},
				thief:    areaID{2, 1},
				opponent: true,
				empty:    []areaID{areaID{3, 6}, areaID{4, 6}},
				expected: []areaID{areaID{2, 6}},
			}),
			Entry("up empty, but thief", entry{
				selected: areaID{6, 5},
				thief:    areaID{2, 5},
				opponent: true,
				empty:    []areaID{areaID{3, 5}, areaID{4, 5}, areaID{5, 5}},
				expected: noAreaIDS,
			}),
		)
	})
})

var _ = Describe("g.camelEEE", func() {
	var (
		c                   *gin.Context
		resp                *httptest.ResponseRecorder
		g                   game
		u1, u2              user.User2
		ta, expectedArea, a area
	)

	type entry struct {
		selected areaID
		thief    areaID
		empty    []areaID
		expected areaID
	}

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()
		g = createGame(c, u1, u2)
	})

	DescribeTable("when valid selected area",
		func(e entry) {
			from := g.grid.area(e.selected.row, e.selected.column)
			ta = g.grid.area(e.thief.row, e.thief.column)
			ta.thief.pid = 1
			g.grid = g.grid.updateArea(ta)

			for _, aid := range e.empty {
				empty := noArea
				empty.areaID = aid
				g.grid = g.grid.updateArea(empty)
			}

			expectedArea = g.grid.area(e.expected.row, e.expected.column)
			a = g.grid.camelEEE(from)
			Expect(a).To(Equal(expectedArea))
		},
		Entry("no area selected", entry{
			selected: noAreaID,
			thief:    areaID{1, 5},
			empty:    []areaID{areaID{2, 4}},
			expected: noAreaID,
		}),
		Entry("thief east east east", entry{
			selected: areaID{1, 1},
			thief:    noAreaID,
			empty:    []areaID{areaID{2, 4}},
			expected: areaID{1, 4},
		}),
		Entry("thief east east east, empty to east", entry{
			selected: areaID{1, 6},
			thief:    areaID{1, 5},
			empty:    []areaID{areaID{1, 3}, areaID{2, 4}},
			expected: noAreaID,
		}),
	)
})

var _ = Describe("g.camelWWW", func() {
	var (
		c                   *gin.Context
		resp                *httptest.ResponseRecorder
		g                   game
		u1, u2              user.User2
		ta, expectedArea, a area
	)

	type entry struct {
		selected areaID
		thief    areaID
		empty    []areaID
		expected areaID
	}

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()
		g = createGame(c, u1, u2)
	})

	DescribeTable("when valid selected area",
		func(e entry) {
			from := g.grid.area(e.selected.row, e.selected.column)
			ta = g.grid.area(e.thief.row, e.thief.column)
			ta.thief.pid = 1
			g.grid = g.grid.updateArea(ta)

			for _, aid := range e.empty {
				empty := noArea
				empty.areaID = aid
				g.grid = g.grid.updateArea(empty)
			}

			expectedArea = g.grid.area(e.expected.row, e.expected.column)
			a = g.grid.camelWWW(from)
			Expect(a).To(Equal(expectedArea))
		},
		Entry("no area selected", entry{
			selected: noAreaID,
			thief:    areaID{1, 5},
			empty:    []areaID{areaID{2, 4}},
			expected: noAreaID,
		}),
		Entry("thief west west west", entry{
			selected: areaID{1, 4},
			thief:    areaID{1, 5},
			empty:    []areaID{areaID{2, 4}},
			expected: areaID{1, 1},
		}),
		Entry("thief west west west, empty to west", entry{
			selected: areaID{1, 4},
			thief:    areaID{1, 5},
			empty:    []areaID{areaID{1, 3}, areaID{2, 4}},
			expected: noAreaID,
		}),
	)
})

var _ = Describe("g.camelNNN", func() {
	var (
		c                   *gin.Context
		resp                *httptest.ResponseRecorder
		g                   game
		u1, u2              user.User2
		ta, expectedArea, a area
	)

	type entry struct {
		selected areaID
		thief    areaID
		empty    []areaID
		expected areaID
	}

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()
		g = createGame(c, u1, u2)
	})

	DescribeTable("when valid selected area",
		func(e entry) {
			from := g.grid.area(e.selected.row, e.selected.column)
			ta = g.grid.area(e.thief.row, e.thief.column)
			ta.thief.pid = 1
			g.grid = g.grid.updateArea(ta)

			for _, aid := range e.empty {
				empty := noArea
				empty.areaID = aid
				g.grid = g.grid.updateArea(empty)
			}

			expectedArea = g.grid.area(e.expected.row, e.expected.column)
			a = g.grid.camelNNN(from)
			Expect(a).To(Equal(expectedArea))
		},
		Entry("no area selected", entry{
			selected: noAreaID,
			thief:    areaID{1, 5},
			empty:    []areaID{areaID{2, 4}},
			expected: noAreaID,
		}),
		Entry("thief north north north", entry{
			selected: areaID{6, 1},
			thief:    areaID{3, 2},
			empty:    []areaID{areaID{6, 2}},
			expected: areaID{3, 1},
		}),
		Entry("thief north north north, empty to north", entry{
			selected: areaID{6, 1},
			thief:    areaID{5, 2},
			empty:    []areaID{areaID{5, 1}, areaID{6, 2}},
			expected: noAreaID,
		}),
	)
})

var _ = Describe("g.grid.turban0Areas", func() {
	var (
		c                 *gin.Context
		resp              *httptest.ResponseRecorder
		g                 game
		u1, u2            user.User2
		expectedAreas, as []area
	)

	BeforeEach(func() {
		resp = httptest.NewRecorder()
		c, _ = gin.CreateTestContext(resp)

		u1, u2 = createUsers()
		g = createGame(c, u1, u2)
	})

	Context("when no area selected", func() {
		It("should return no areas", func() {
			Expect(g.grid.carpetAreas(noArea)).To(HaveLen(0))
		})
	})

	Context("when full grid without any thieves", func() {
		var ta area

		type entry struct {
			selected areaID
			thief    areaID
			empty    []areaID
			expected []areaID
		}

		DescribeTable("when valid selected area",
			func(e entry) {
				from := g.grid.area(e.selected.row, e.selected.column)
				from.thief.pid = 1
				g.grid = g.grid.updateArea(from)

				if e.thief != noAreaID {
					ta = g.grid.area(e.thief.row, e.thief.column)
					ta.thief.pid = 1
					g.grid = g.grid.updateArea(ta)
				}

				for _, aid := range e.empty {
					a := g.grid.area(aid.row, aid.column)
					a = a.empty()
					g.grid = g.grid.updateArea(a)
				}

				expectedAreas = make([]area, len(e.expected))
				for i, aid := range e.expected {
					expectedAreas[i] = g.grid.area(aid.row, aid.column)
				}
				as = g.grid.turban0Areas(from)
				Expect(as).To(ConsistOf(expectedAreas))
			},
			Entry("east thief", entry{
				selected: areaID{1, 1},
				thief:    areaID{1, 2},
				empty:    noAreaIDS,
				expected: []areaID{areaID{2, 1}},
			}),
			Entry("east empty", entry{
				selected: areaID{1, 1},
				thief:    noAreaID,
				empty:    []areaID{areaID{1, 2}},
				expected: []areaID{areaID{2, 1}},
			}),
			Entry("east, but cannot take second step", entry{
				selected: areaID{1, 1},
				thief:    noAreaID,
				empty:    []areaID{areaID{2, 1}, areaID{2, 2}, areaID{1, 3}},
				expected: noAreaIDS,
			}),
			Entry("west thief", entry{
				selected: areaID{2, 3},
				thief:    areaID{2, 2},
				empty:    []areaID{areaID{2, 4}, areaID{1, 3}, areaID{3, 3}},
				expected: noAreaIDS,
			}),
			Entry("west empty", entry{
				selected: areaID{4, 8},
				thief:    noAreaID,
				empty:    []areaID{areaID{4, 7}, areaID{2, 8}, areaID{3, 7}},
				expected: []areaID{areaID{5, 8}},
			}),
			Entry("west, but cannot take second step", entry{
				selected: areaID{3, 7},
				thief:    noAreaID,
				empty: []areaID{areaID{3, 8}, areaID{2, 7}, areaID{4, 7}, areaID{2, 6}, areaID{4, 6},
					areaID{3, 5}},
				expected: noAreaIDS,
			}),
			Entry("south thief", entry{
				selected: areaID{4, 4},
				thief:    areaID{5, 4},
				empty:    []areaID{areaID{3, 4}, areaID{4, 3}, areaID{4, 5}},
				expected: noAreaIDS,
			}),
			Entry("south empty", entry{
				selected: areaID{5, 8},
				thief:    noAreaID,
				empty:    []areaID{areaID{6, 8}},
				expected: []areaID{areaID{4, 8}, areaID{5, 7}},
			}),
			Entry("south, but cannot take second step", entry{
				selected: areaID{3, 7},
				thief:    noAreaID,
				empty: []areaID{areaID{2, 7}, areaID{3, 6}, areaID{3, 8}, areaID{5, 7}, areaID{4, 6},
					areaID{4, 8}},
				expected: noAreaIDS,
			}),
			Entry("north thief", entry{
				selected: areaID{5, 4},
				thief:    areaID{4, 4},
				empty:    []areaID{areaID{6, 4}, areaID{5, 3}, areaID{5, 5}},
				expected: noAreaIDS,
			}),
			Entry("north empty", entry{
				selected: areaID{5, 8},
				thief:    noAreaID,
				empty:    []areaID{areaID{4, 8}},
				expected: []areaID{areaID{6, 8}, areaID{5, 7}},
			}),
			Entry("north, but cannot take second step", entry{
				selected: areaID{3, 7},
				thief:    noAreaID,
				empty: []areaID{areaID{4, 7}, areaID{3, 6}, areaID{3, 8}, areaID{1, 7}, areaID{2, 6},
					areaID{2, 8}},
				expected: noAreaIDS,
			}),
		)
	})
})
